# A database project
resource "nuodbaas_database" "nuodb" {
  organization = nuodbaas_project.nuodb.organization
  project      = nuodbaas_project.nuodb.name
  name         = "nuodb"
  dba_password = "helloworld"
}

# A database with more fields set
resource "nuodbaas_database" "dbaas" {
  organization = nuodbaas_project.nuodb.organization
  project      = nuodbaas_project.nuodb.name
  name         = "dbaas"
  tier         = "n0.nano"
  dba_password = "helloworld"
  maintenance = {
    is_disabled = false
  }

  properties = {
    tier_parameters = {
      zones        = jsonencode(["us-east-2a", "us-east-2c"])
      capacityType = "spot"
    }
  }

  timeouts {
    create = "1m"
    update = "1m"
  }
}