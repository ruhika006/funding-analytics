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
    docker run -d --name myapp -p 8080:8080 ruhika0817/db-app
  EOF


  tags = {
    Name = "ruh-ubuntu-ec2"
  }

}

