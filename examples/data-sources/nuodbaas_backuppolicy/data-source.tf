# Data source that returns the attributes of a specific backup policy
data "nuodbaas_backuppolicy" "policy_details" {
  organization = "org"
  name         = "pol"
  depends_on   = [nuodbaas_backuppolicy.pol]
}
