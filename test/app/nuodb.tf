resource "nuodbaas_project" "test" {
  organization = var.org_name
  name         = var.project_name
  sla          = "test"
  tier         = "n0.nano"
}

# Add a database into the project
resource "nuodbaas_database" "app" {
  organization = nuodbaas_project.test.organization
  project      = nuodbaas_project.test.name
  name         = var.db_name
  dba_password = sensitive(var.dba_password)
}
