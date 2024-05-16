# A backup policy with weekly backup frequency on all databases in organization
resource "nuodbaas_backuppolicy" "basic" {
  organization = "org"
  name         = "basic"
  frequency    = "@weekly"

  selector = {
    scope = "org"
  }
}

# A backup policy with daily backup frequency on databases with SLA prod and label
resource "nuodbaas_backuppolicy" "daily" {
  organization = "org"
  name         = "daily"
  frequency    = "@daily"

  selector = {
    scope = "org"
    slas  = ["prod"]

    labels = {
      "rpo" : "1d"
    }
  }
}
