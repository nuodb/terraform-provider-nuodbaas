# Get all databases
data "nuodbaas_databases" "databaseList" {}


# Get all databases in a given organization
data "nuodbaas_databases" "databaseList" {
  filter {
    organization = "system"
  }
}

# Get all databases in a given project
data "nuodbaas_databases" "databaseList" {
  filter {
    organization = "system"
    project      = "nuodb"
  }
}