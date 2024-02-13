# Get details about a single database
data "nuodbaas_database" "database_details" {
  name         = "dbaas"
  organization = "system"
  project      = "nuodb"
}