resource "terraform_data" "dbapp" {

  provisioner "local-exec" {
    command = "./db_connect.sh"

    environment = {
      DB_NAME      = nuodbaas_database.app.name
      DB_USER      = "dba"
      DB_PASSWORD  = var.dba_password
      PEER_ADDRESS = "${nuodbaas_database.app.status.sql_endpoint}:443"
      CA_CERT      = nuodbaas_database.app.status.ca_pem
    }
  }
}
