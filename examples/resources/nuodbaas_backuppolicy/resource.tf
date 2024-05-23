# A backup policy with minimal configuration
resource "nuodbaas_backuppolicy" "basic" {
  organization = "org"
  name         = "basic"
  frequency    = "@weekly"
  selector = {
    scope = "org"
  }
}

# A backup policy with explicit configuration for various attributes
resource "nuodbaas_backuppolicy" "pol" {
  organization = "org"
  name         = "pol"
  labels = {
    "provisioned-by" : "terraform"
  }
  frequency = "@daily"
  selector = {
    scope = "org"
    slas  = ["qa", "prod"]
    tiers = ["n0.small", "n1.small"]
    labels = {
      "rpo" : "1d"
    }
  }
  retention = {
    hourly  = 24
    daily   = 7
    weekly  = 4
    monthly = 12
    yearly  = 3
  }
  suspended = false
  properties = {
    propagate_policy_labels   = true
    propagate_database_labels = true
  }
}
