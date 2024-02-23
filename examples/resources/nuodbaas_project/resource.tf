# A project with minimal configuration
resource "nuodbaas_project" "basic" {
  organization = "org"
  name         = "basic"
  sla          = "dev"
  tier         = "n0.nano"
}

# A project with explicit configuration for various attributes
resource "nuodbaas_project" "proj" {
  organization = "org"
  name         = "proj"
  sla          = "prod"
  tier         = "n0.nano"
  labels = {
    color  = "blue"
    flavor = "mild"
  }
  properties = {
    product_version = "5.1"
    tier_parameters = {
      zone  = "us-east"
      group = "dev"
    }
  }
}
