# Data source that returns the fully-qualified names of all backup policies
data "nuodbaas_backuppolicies" "policy_list" {}

# Data source that returns the fully-qualified names of backup policies within an organization
data "nuodbaas_backuppolicies" "org_policy_list" {
  filter = {
    organization = "org"
  }
}

# Data source that returns the fully-qualified names of backup policies satisfying label requirements
data "nuodbaas_backuppolicies" "label_policy_list" {
  filter = {
    labels = ["withkey", "key=expected", "key!=unexpected", "!withoutkey"]
  }
}
