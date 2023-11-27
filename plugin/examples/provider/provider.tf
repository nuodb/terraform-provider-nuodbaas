terraform {
  required_providers {
    nuodbaas = {
      source = "hashicorp.com/edu/nuodbaas"
    }
  }
}

provider "nuodbaas" {
  organization=var.dbaas_credentials.organization
  username=var.dbaas_credentials.username
  password=sensitive(var.dbaas_credentials.password)
  url_base="http://localhost/nuodb-cp"
}

resource "nuodbaas_project" "dev"{
  organization=var.dbaas_credentials.organization
  name="dev"
  sla="dev"
  tier="n1.small"
  maintenance = {
    is_disabled=false
  }
}

resource "nuodbaas_project" "nuodb" {
  organization=var.dbaas_credentials.organization
  name="nuodb"
  sla="dev"
  tier="n1.small"
}

resource "nuodbaas_database" "nuodb" { 
  organization=var.dbaas_credentials.organization
  project=nuodbaas_project.nuodb.name
  name="nuodb"
  tier="n0.nano"
  dba_password="helloworld"
  maintenance = {
    expires_in="2d"
  }

  properties = {
    archive_disk_size = "20Gi"
  }
}

data "nuodbaas_projects" "projectsList" {
  organization=var.dbaas_credentials.organization
}

output "projectsList" {
  value = data.nuodbaas_projects.projectsList
}

