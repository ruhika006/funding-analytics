terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region  = "ap-south-2"
  profile = "practice"
}


# For ECS

/* Uncomment below code block to use ECS Fargate

data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

resource "aws_iam_role" "ecs_task_execution" {
  name = "ecs-task-execution"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
    }]
  })
}

resource "aws_security_group" "ecs_sg" {
  name = "ecs-sg"
  vpc_id = data.aws_vpc.default.id

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_ecs_cluster" "main" {
  name = "db-app-cluster"
}

resource "aws_ecs_task_definition" "main" {
  family                   = "db-app"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  container_definitions = jsonencode([{
    name  = "db-app"
    image = "ruhika0817/db-app"
  }])
}

resource "aws_ecs_service" "main" {
  name            = "db-app-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.main.arn
  desired_count   = 1
  launch_type     = "FARGATE"
  network_configuration {
    subnets            = data.aws_subnets.default.ids
    security_groups    = [aws_security_group.ecs_sg.id]
    assign_public_ip   = true
  }
}

*/

# EC2

/* Uncommnet below code block to use AppRunner

resource "aws_instance" "ruh-ubuntu-ec2"{
  ami           = "ami-0e7938ad51d883574"
  instance_type = "t3.micro"
  key_name      = "ruhika-key"

  user_data = <<-EOF
    #!/bin/bash

    # Update system
    apt-get update -y

    # Install Docker
    apt-get install -y docker.io
    systemctl start docker
    systemctl enable docker

    # Pull your Docker image
    docker pull ruhika0817/db-app

    # Run container
    docker run -d --name myapp -p 80:80 ruhika0817/db-app
  EOF


  tags = {
    Name = "ruh-ubuntu-ec2"
  }

}
*/



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

