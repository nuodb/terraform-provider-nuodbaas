# Get all projects
data "nuodbaas_projects" "projectsList" {}

# Get all projects in a given organization
data "nuodbaas_projects" "projectsList" {
  filter {
    organization = "system"
  }
}