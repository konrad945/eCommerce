terraform {
  required_providers {
    kind = {
      source = "tehcyx/kind"
      version = "0.0.13"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "2.19.0"
    }
  }

  required_version = ">= 1.2.5"
}