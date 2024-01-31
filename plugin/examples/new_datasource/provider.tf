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
  base_url     = "https://harsh.us-east-2.mdsol.mynuodb.com/api"
}

resource "nuodbaas_project" "nuodb" {
  organization = var.dbaas_credentials.organization
  name         = "nuodb"
  sla          = "dev"
  tier         = "nx.small"
  maintenance = {
    expires_in = "5d"
  }

  # properties = {}
}


resource "nuodbaas_database" "dbaas" {
  organization = var.dbaas_credentials.organization
  project      = nuodbaas_project.nuodb.name
  name         = "dbaas"
  tier         = "nx.small"
  dba_password = "helloworld"
  maintenance = {
    expires_in = "2d"
  }

  properties = {
    archive_disk_size = "1Gi"
    tier_parameters = {
      smReplicas = 1
      teReplicas = 2
      zones = jsonencode([ "us-east-2a", "us-east-2c" ])
      capacityType = jsonencode("spot")
    }
  }
}

# data "nuodbaas_projects" "projectsList" {
#   filter {
#     organization = var.dbaas_credentials.organization
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

# data "nuodbaas_databases" "databaseList" {
#   filter {
#     organization = var.dbaas_credentials.organization
#     project      = "nuodb"
#   }
# }

# output "databaseList" {
#   value = data.nuodbaas_databases.databaseList
# }
# locals {
#   database_names = {
#     for database in data.nuodbaas_databases.databaseList.databases :
#     database.name => database
#   }
# }

# data "nuodbaas_database" "databaseDetail" {
#   for_each = local.database_names
#   organization = var.dbaas_credentials.organization
#   project      = nuodbaas_project.nuodb.name
#   name = "${each.key}"
# }

# output "nuodbaas_databaseDetails" {
#   value = [for database in data.nuodbaas_database.databaseDetail: database]
# }

