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
  tier="n0.nano"
  maintenance = {
    expires_in="1d"
  }
}

resource "nuodbaas_project" "nuodb" {
  organization=var.dbaas_credentials.organization
  name="dassault"
  sla="dev"
  tier="n1.small"

  maintenance = {
    expires_in="5d"
  }
}
