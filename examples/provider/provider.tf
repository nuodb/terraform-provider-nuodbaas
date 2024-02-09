terraform {
  required_providers {
    nuodbaas = {
      source  = "registry.terraform.io/nuodb/nuodbaas"
      version = "0.1.0"
    }
  }
}

provider "nuodbaas" {
  user     = var.dbaas_credentials.user
  password = sensitive(var.dbaas_credentials.password)
  url_base = var.dbaas_credentials.url_base
}

# Create a basic project
resource "nuodbaas_project" "proj" {
  organization = "org"
  name         = "nuodb"
  sla          = "prod"
  tier         = "n0.nano"
}

# Add a database into the project
resource "nuodbaas_database" "db" {
  organization = nuodbaas_project.proj.organization
  project      = nuodbaas_project.proj.name
  name         = "db"
  dba_password = "helloworld"

  # By using the attributes of nuodbaas_project.proj for organization and project,
  # we are already creating an implicit dependency on the project being created
  # before we add a database to it. If you do not want to (or cannot) depend on
  # implicit dependencies, you can create an explicit one:
  depends_on = [nuodbaas_project.proj]
}

data "nuodbaas_database" "database_details" {
  name         = nuodbaas_database.db.name
  organization = nuodbaas_database.db.organization
  project      = nuodbaas_database.db.project

  # By using the attributes of nuodbaas_database.db as arguments, we are
  # creating an implicit dependency on the database being created before we
  # try to fetch it. If you do not want to (or cannot) depend on
  # implicit dependencies, you can create an explicit one:
  depends_on = [nuodbaas_database.db]
}