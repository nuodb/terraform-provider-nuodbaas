variable "org_name" {
  description = "Organization name"
  type        = string
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "testproj"
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "app"
}

variable "dba_password" {
  description = "Password for the initial dba user"
  type        = string
  default     = "password"
}