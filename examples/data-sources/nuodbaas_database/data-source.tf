# Data source that returns the attributes of a specific database
data "nuodbaas_database" "database_details" {
  organization = "org"
  project      = "proj"
  name         = "db"
  depends_on   = [nuodbaas_database.db]
}
