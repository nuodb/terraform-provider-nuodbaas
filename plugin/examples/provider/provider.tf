terraform {
  required_providers {
    nuodbaas = {
      source = "hashicorp.com/edu/nuodbaas"
    }
  }
}

provider "nuodbaas" {
  organization=var.dbaas_credentials.organization
  username=var.dbaas_credentials.username
  password=sensitive(var.dbaas_credentials.password)
  url_base="http://localhost/nuodb-cp"
}

resource "nuodbaas_project" "dev"{
  organization=var.dbaas_credentials.organization
  name="dev"
  sla="dev"
  tier="n1.small"
  maintenance = {
    is_disabled=false
  }

  properties = {}
}

resource "nuodbaas_project" "nuodb" {
  organization=var.dbaas_credentials.organization
  name="nuodb"
  sla="dev"
  tier="n1.small"
}

resource "nuodbaas_database" "nuodb" {
  organization = var.dbaas_credentials.organization
  project      = nuodbaas_project.nuodb.name
  name         = "nuodb"
  tier         = "n0.nano"
  dba_password = "helloworld"
  properties = {
    archive_disk_size = "7Gi"
  }

}

resource "nuodbaas_database" "dbaas" {
  organization = var.dbaas_credentials.organization
  project      = nuodbaas_project.nuodb.name
  name         = "dbaas"
  # tier         = "n1.small"
  dba_password = "helloworld"
  maintenance = {
    expires_in = "2d"
  }

  # properties = {
  #   archive_disk_size = "1Gi"
  # }
}

# data "nuodbaas_projects" "projectsList" {
#   filter {
    
#   }
# }

# output "projectsList" {
#   value = data.nuodbaas_projects.projectsList
# }

# locals {
#   project_names = {
#     for project in data.nuodbaas_projects.projectsList.projects :
#     project.name => project
#   }
# }

# data "nuodbaas_project" "projectDetail" {
#   for_each = local.project_names
#   organization = var.dbaas_credentials.organization
#   name = "${each.key}"
# }

# output "projectDetail" {
#   value = data.nuodbaas_project.projectDetail
# }


