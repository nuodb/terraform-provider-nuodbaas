resource "docker_container" "dbapp" {
  name  = "dbapp"
  image = "python:3.12.2-alpine3.19"
  env = [ 
    "DB_NAME=${nuodbaas_database.app.name}",
    "DB_USER=dba",
    "DB_PASSWORD=${var.dba_password}",
    "PEER_ADDRESS=${nuodbaas_database.app.status.sql_endpoint}:443",
    "CA_CERT=${nuodbaas_database.app.status.ca_pem}"
    ]
  mounts {
    type = "bind"
    target = "/usr/src/dbapp"
    source = "${path.cwd}"
  }
  command = [ "/usr/src/dbapp/db_connect.sh" ]
  must_run = false
  attach = true
  logs = true
}

output "exit" {
  value = docker_container.dbapp.exit_code
}

output "logs" {
  value = docker_container.dbapp.container_logs
}