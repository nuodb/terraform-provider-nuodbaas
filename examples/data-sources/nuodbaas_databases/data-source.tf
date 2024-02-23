# Data source that returns the fully-qualified names of all databases
data "nuodbaas_databases" "database_list" {}


# Data source that returns the fully-qualified names of databases within an organization
data "nuodbaas_databases" "org_database_list" {
  filter = {
    organization = "org"
  }
}

# Data source that returns the fully-qualified names of databases within a project
data "nuodbaas_databases" "proj_database_list" {
  filter = {
    organization = "org"
    project      = "proj"
  }
}

# Data source that returns the fully-qualified names of databases satisfying label requirements
data "nuodbaas_databases" "label_database_list" {
  filter = {
    labels = ["withkey", "key=expected", "key!=unexpected", "!withoutkey"]
  }
}
