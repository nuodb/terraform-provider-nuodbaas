terraform {
  required_providers {
    nuodbaas = {
      source = "registry.terraform.io/nuodb/nuodbaas"
    }
  }
}

provider "nuodbaas" {
  organization = var.dbaas_credentials.organization
  username     = var.dbaas_credentials.username
  password     = sensitive(var.dbaas_credentials.password)
  url_base     = "http://localhost/nuodb-cp"
}

# Create a basic project
resource "nuodbaas_project" "nuodb" {
  organization = var.dbaas_credentials.organization
  name         = "nuodb"
  sla          = "prod"
  tier         = "n0.nano"
}

# Add a database into the project
resource "nuodbaas_database" "nuodb" {
  organization = nuodbaas_project.nuodb.organization
  project      = nuodbaas_project.nuodb.name
  name         = "nuodb"
  dba_password = "helloworld"
}