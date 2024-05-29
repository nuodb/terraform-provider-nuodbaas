# Data source that returns the attributes of a specific database
data "nuodbaas_database" "database_details" {
  organization = nuodbaas_database.db.organization
  project      = nuodbaas_database.db.project
  name         = nuodbaas_database.db.name
}
