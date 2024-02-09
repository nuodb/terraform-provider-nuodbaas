# Get all projects
data "nuodbaas_projects" "projects_list" {}

# Get all projects in a given organization
data "nuodbaas_projects" "org_projects_list" {
  filter {
    organization = "system"
  }
}