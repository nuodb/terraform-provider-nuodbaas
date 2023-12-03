variable "dbaas_credentials" {
  type = object({
    username     = string
    organization = string
    password     = string
  })

  default = {
    organization = "system"
    username     = "admin"
    password     = "Iw(S0>=X?t"
  }
}

# Iw(S0>=X?t