provider "aws" {
  region  = "ap-south-1"
  profile = "practice"
  alias = "apprunner"
}


resource "aws_iam_role" "apprunner_ecr_role" {
  name = "apprunner-ecr-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = { Service = "build.apprunner.amazonaws.com" }
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "apprunner_ecr_policy" {
  role       = aws_iam_role.apprunner_ecr_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSAppRunnerServicePolicyForECRAccess"
}


resource "aws_apprunner_service" "db_apprunner" {
  service_name = "db_apprunner"
  provider     = aws.apprunner

  source_configuration {
    authentication_configuration {
      access_role_arn = aws_iam_role.apprunner_ecr_role.arn
    }

    image_repository {
      image_identifier      = "912390896173.dkr.ecr.ap-south-1.amazonaws.com/db/db-app:latest"
      image_repository_type = "ECR"

      image_configuration {
        port = "8080"
      }
    }

    auto_deployments_enabled = false
  }

  tags = {
    Name = "db-apprunner-service"
  }
}