# Data source that returns the attributes of a specific project
data "nuodbaas_project" "project_details" {
  organization = nuodbaas_project.proj.organization
  name         = nuodbaas_project.proj.name
}
