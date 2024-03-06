terraform {
  required_providers {
    nuodbaas = {
      source  = "registry.terraform.io/nuodb/nuodbaas"
      version = "1.0.0"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

provider "docker" { }

provider "nuodbaas" {
  # Don't wait for the project to be created to proceed with creating other resources.
  # DBaaS will block the database startup until project is ready.
  timeouts = {
    project = {
      create = "0"
    }
  }
 }
