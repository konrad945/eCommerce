terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = "2.12.1"
    }
    helm = {
      source = "hashicorp/helm"
      version = "2.6.0"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "2.19.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "3.1.1"
    }
  }

  required_version = ">= 1.2.5"
}