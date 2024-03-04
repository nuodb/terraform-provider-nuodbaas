terraform {
  required_providers {
    nuodbaas = {
      source  = "registry.terraform.io/nuodb/nuodbaas"
      version = "0.2.0"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

provider "docker" { }

provider "nuodbaas" { }
