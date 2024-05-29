# A backup referencing an existing backup handle
resource "nuodbaas_backup" "backup" {
  organization = nuodbaas_database.db.organization
  project      = nuodbaas_database.db.project
  database     = nuodbaas_database.db.name
  name         = "backup"
  labels = {
    latest  = "true"
    purpose = "test"
  }
  import_source = {
    backup_handle = "backup"
    backup_plugin = "embedded.cp.nuodb.com"
  }
}
