terraform {
  required_providers {
    nuodbaas = {
      source  = "registry.terraform.io/nuodb/nuodbaas"
      version = "0.2.0"
    }
  }
}

provider "nuodbaas" {
  user     = var.dbaas_credentials.user
  password = sensitive(var.dbaas_credentials.password)
  url_base = var.dbaas_credentials.url_base
  timeouts = {
    defaults = {
      create = "10m"
      update = "5m"
      delete = "30s"
    }
  }
}

# Create a project
resource "nuodbaas_project" "proj" {
  organization = "org"
  name         = "proj"
  sla          = "dev"
  tier         = "n0.nano"
}

# Create a database within the project
resource "nuodbaas_database" "db" {
  organization = nuodbaas_project.proj.organization
  project      = nuodbaas_project.proj.name
  name         = "db"
  dba_password = "secret"

  # By using the attributes of nuodbaas_project.proj for organization and
  # project, we are defining an implicit dependency on the project, which
  # causes it to be created before the database. An explicit dependency can be
  # defined as follows:
  depends_on = [nuodbaas_project.proj]
}

data "nuodbaas_database" "database_details" {
  organization = nuodbaas_database.db.organization
  project      = nuodbaas_database.db.project
  name         = nuodbaas_database.db.name

  # By using the attributes of nuodbaas_database.db, we are defining an
  # implicit dependency on the database, which causes it to be created before
  # the data source is read. An explicit dependency can be defined as follows:
  depends_on = [nuodbaas_database.db]
}
