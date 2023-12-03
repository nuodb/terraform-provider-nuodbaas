terraform {
  required_providers {
    nuodbaas = {
      source = "hashicorp.com/edu/nuodbaas"
    }
  }
}

provider "nuodbaas" {
  organization = var.dbaas_credentials.organization
  username     = var.dbaas_credentials.username
  password     = sensitive(var.dbaas_credentials.password)
  base_url     = "https://harsh.us-east-2.mdsol.mynuodb.com/api"
}

# resource "nuodbaas_project" "nuodb" {
#   organization = var.dbaas_credentials.organization
#   name         = "nuodb"
#   sla          = "dev"
#   tier         = "nx.small"
#   maintenance = {
#     expires_in = "5d"
#   }
# }


# resource "nuodbaas_database" "dbaas" {
#   organization = var.dbaas_credentials.organization
#   project      = nuodbaas_project.nuodb.name
#   name         = "dbaas"
#   tier         = "nx.small"
#   dba_password = "helloworld"
#   maintenance = {
#     expires_in = "2d"
#   }

#   properties = {
#     archive_disk_size = "1Gi"
#   }
# }

data "nuodbaas_projects" "projectsList" {
  filter = {
    organization = var.dbaas_credentials.organization
  }
}

output "projectsList" {
  value = data.nuodbaas_projects.projectsList
}

# data "nuodbaas_databases" "databaseList" {
#   filter = {
#     organization = var.dbaas_credentials.organization
#     project      = nuodbaas_project.nuodb.name
#   }
# }

# output "databaseList" {
#   value = data.nuodbaas_databases.databaseList
# }

# data "nuodbaas_database" "databaseDetail" {
#   depends_on = []
#   for_each = toset(data.nuodbaas_databases.databaseList.databases == null ? [] : data.nuodbaas_databases.databaseList.databases)
#   filter = {
#     organization = var.dbaas_credentials.organization
#     project      = nuodbaas_project.nuodb.name
#   }
#   name = "${each.key}"
# }

# output "nuodbaas_databaseDetails" {
#   value = [for database in data.nuodbaas_database.databaseDetail: database.database]
# }

