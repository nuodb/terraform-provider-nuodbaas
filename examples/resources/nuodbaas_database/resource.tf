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
    archive_disk_size = "10Gi"
    tier_parameters = {
      zones        = jsonencode(["us-east-2a", "us-east-2c"])
      capacityType = "spot"
    }
  }
}
