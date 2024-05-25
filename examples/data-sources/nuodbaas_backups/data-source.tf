# Data source that returns the fully-qualified names of all backups
data "nuodbaas_backups" "backup_list" {}


# Data source that returns the fully-qualified names of backups within an organization
data "nuodbaas_backups" "org_backup_list" {
  filter = {
    organization = "org"
  }
}

# Data source that returns the fully-qualified names of backups within a project
data "nuodbaas_backups" "proj_backup_list" {
  filter = {
    organization = "org"
    project      = "proj"
  }
}

# Data source that returns the fully-qualified names of backups within a database
data "nuodbaas_backups" "db_backup_list" {
  filter = {
    organization = "org"
    project      = "proj"
    database     = "db"
  }
}

# Data source that returns the fully-qualified names of backups satisfying label requirements
data "nuodbaas_backups" "label_backup_list" {
  filter = {
    labels = ["withkey", "key=expected", "key!=unexpected", "!withoutkey"]
  }
}
