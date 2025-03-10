terraform {
  required_providers {
    nuodbaas = {
      source = "registry.terraform.io/nuodb/nuodbaas"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.6.0"
    }
  }
}

# Used to generate unique project name to avoid collisions
provider "random" {}

provider "nuodbaas" {
  # Don't wait for the project to be created to proceed with creating other resources.
  # DBaaS will block the database startup until project is ready.
  timeouts = {
    project = {
      create = "0"
    }
  }
}
