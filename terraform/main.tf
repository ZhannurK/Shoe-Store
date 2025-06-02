terraform {
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
  }

  required_version = ">= 1.0.0"
}

resource "null_resource" "docker_compose_up" {
  provisioner "local-exec" {
    command = "docker compose -f ../docker-compose.yml up -d"
  }
}

resource "null_resource" "docker_compose_down" {
  provisioner "local-exec" {
    when    = destroy
    command = "docker compose -f ../docker-compose.yml down"
  }
}
