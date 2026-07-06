# A database with minimal configuration
resource "nuodbaas_database" "basic" {
  organization = nuodbaas_project.proj.organization
  project      = nuodbaas_project.proj.name
  name         = "basic"
  dba_password = "secret"
}

# A database with explicit configuration for various attributes
resource "nuodbaas_database" "db" {
  organization = nuodbaas_project.proj.organization
  project      = nuodbaas_project.proj.name
  name         = "db"
  tier         = "n0.nano"
  dba_password = "secret"
  labels = {
    color  = "green"
    flavor = "bold"
  }
  properties = {
    journal_disk_size = "1Gi"
    tier_parameters = {
      zones        = jsonencode(["us-east-2a", "us-east-2c"])
      capacityType = "spot"
    }
    archive_disk_auto_resize = {
      initial_size = "20Gi"
      max_size     = "500Gi"
      threshold = {
        percentage_available = 5
      }
    }
  }
}
