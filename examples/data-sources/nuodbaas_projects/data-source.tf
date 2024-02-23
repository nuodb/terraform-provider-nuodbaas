# Data source that returns the fully-qualified names of all projects
data "nuodbaas_projects" "project_list" {}

# Data source that returns the fully-qualified names of projects within an organization
data "nuodbaas_projects" "org_project_list" {
  filter = {
    organization = "org"
  }
}

# Data source that returns the fully-qualified names of projects satisfying label requirements
data "nuodbaas_projects" "label_project_list" {
  filter = {
    labels = ["withkey", "key=expected", "key!=unexpected", "!withoutkey"]
  }
}
