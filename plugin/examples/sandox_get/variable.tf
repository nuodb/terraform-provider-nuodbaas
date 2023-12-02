variable "dbaas_credentials" {
  type = object({
    username     = string
    organization = string
    password     = string
  })

  default = {
    organization = "system"
    username     = "admin"
    password     = "fJ_Mck1RsR"
  }
  
}
