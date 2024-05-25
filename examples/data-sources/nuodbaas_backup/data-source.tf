# Data source that returns the attributes of a specific backup
data "nuodbaas_backup" "backup_details" {
  organization = nuodbaas_backup.backup.organization
  project      = nuodbaas_backup.backup.project
  database     = nuodbaas_backup.backup.database
  name         = nuodbaas_backup.backup.name
}
