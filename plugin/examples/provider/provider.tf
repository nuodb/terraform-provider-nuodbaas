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

resource "nuodbaas_project" "nuodb" {
  organization=var.dbaas_credentials.organization
  name="nuodb"
  sla="dev"
  tier="n0.small"

  maintenance = {
    expires_in="5d"
  }
}

resource "nuodbaas_database" "nuodb" {
  organization=var.dbaas_credentials.organization
  project=nuodbaas_project.nuodb.name
  name="nuodb"
  tier="n1.small"
  password="helloworld"
  maintenance = {
    expires_in="2d"
  }

  # archive_disk_size = "15Gi"
  # journal_disk_size = "10Gi"

  # properties = {}
}
