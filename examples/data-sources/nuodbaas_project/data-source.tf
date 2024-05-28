# Data source that returns the attributes of a specific project
data "nuodbaas_project" "project_details" {
  organization = "org"
  name         = "proj"
  depends_on   = [nuodbaas_project.proj]
}
