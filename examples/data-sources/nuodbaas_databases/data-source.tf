# Get all databases
data "nuodbaas_databases" "database_list" {}


# Get all databases in a given organization
data "nuodbaas_databases" "org_database_list" {
  filter = {
    organization = "system"
  }
}

# Get all databases in a given project
data "nuodbaas_databases" "proj_database_list" {
  filter = {
    organization = "system"
    project      = "nuodb"
  }
}
