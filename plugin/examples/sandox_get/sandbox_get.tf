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
  base_url     = "https://dbaasuisbx.us-east-1.dbaas.nuodb.io/api"
}


data "nuodbaas_projects" "projectsList" {
  filter = {
    organization = var.dbaas_credentials.organization
  }
  
}

output "projectsList" {
  value = data.nuodbaas_projects.projectsList
}

data "nuodbaas_databases" "databaseList" {
  filter = {
    organization = var.dbaas_credentials.organization
    project      = "support"
  }
}

output "databaseList" {
  value = data.nuodbaas_databases.databaseList
}

data "nuodbaas_database" "databaseDetail" {
  depends_on = []
  for_each = toset(data.nuodbaas_databases.databaseList.databases == null ? [] : data.nuodbaas_databases.databaseList.databases)
  filter = {
    organization = var.dbaas_credentials.organization
    project      = "support"
  }
  name = "${each.key}"
}

output "nuodbaas_databaseDetails" {
  value = [for database in data.nuodbaas_database.databaseDetail: database.database]
}

