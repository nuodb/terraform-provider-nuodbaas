provider "nuodbaas" {
  # Credentials and URL supplied by environment variables

  # Allow project and database to become ready in parallel
  timeouts = {
    project = {
      create = "0"
    }
  }
}

# Create a project
resource "nuodbaas_project" "proj" {
  organization = "acme"
  name         = "messaging"
  sla          = "dev"
  tier         = "n0.nano"
}

# Create a database within the project
resource "nuodbaas_database" "db" {
  organization = nuodbaas_project.proj.organization
  project      = nuodbaas_project.proj.name
  name         = "demo"
  dba_password = "changeIt"
}

# Expose nuosql arguments to connect to database
output "nuosql_args" {
  value     = <<-EOT
  ${nuodbaas_database.db.name}@${nuodbaas_database.db.status.sql_endpoint}:443 \
      --user dba --password ${nuodbaas_database.db.dba_password} \
      --connection-property trustedCertificates='${nuodbaas_database.db.status.ca_pem}'
  EOT
  sensitive = true
}
