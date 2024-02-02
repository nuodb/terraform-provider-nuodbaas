# Get details about a single project
data "nuodbaas_project" "projectDetails" {
  name         = "nuodb"
  organization = "system"
}