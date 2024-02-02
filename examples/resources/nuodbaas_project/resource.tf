# A basic project
resource "nuodbaas_project" "nuodb" {
  organization = var.dbaas_credentials.organization
  name         = "nuodb"
  sla          = "prod"
  tier         = "n0.nano"
}

# A project with more fields set
resource "nuodbaas_project" "dev" {
  organization = var.dbaas_credentials.organization
  name         = "dev"
  sla          = "dev"
  tier         = "n0.nano"
  maintenance = {
    is_disabled = false
  }

  properties = {
    tier_parameters = {
      zone  = "us-east"
      group = "dev"
    }
  }

  timeouts {
    create = "5m"
    update = "5m"
  }
}