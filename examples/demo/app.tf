provider "docker" {}

resource "docker_image" "python" {
  name = "python:3.12.2-alpine3.19"
}

resource "docker_container" "dbapp" {
  name  = "dbapp"
  image = docker_image.python.image_id
  env = [
    "DB_NAME=${nuodbaas_database.db.name}",
    "DB_USER=dba",
    "DB_PASSWORD=${nuodbaas_database.db.dba_password}",
    "DB_HOST=${nuodbaas_database.db.status.sql_endpoint}:443",
    "CA_CERT=${nuodbaas_database.db.status.ca_pem}"
  ]
  mounts {
    type   = "bind"
    target = "/dbapp"
    source = "${path.cwd}/dbapp"
  }
  command = ["/dbapp/app.sh"]
}
