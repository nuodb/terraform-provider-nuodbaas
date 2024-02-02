# Get details about a single database
data "nuodbaas_database" "databaseDetails" {
  name         = "dbaas"
  organization = "system"
  project      = "nuodb"
}