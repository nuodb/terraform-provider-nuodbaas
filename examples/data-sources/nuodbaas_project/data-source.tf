# Get details about a single project
data "nuodbaas_project" "project_details" {
  name         = "nuodb"
  organization = "system"
}