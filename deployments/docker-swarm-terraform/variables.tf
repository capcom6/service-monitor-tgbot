variable "swarm-manager-host" {
  type        = string
  sensitive   = true
  description = "Address of swarm manager"
}

variable "app-name" {
  type        = string
  description = "Name of app"
}

variable "app-version" {
  type        = string
  description = "Version of Docker image of app"
  default     = "master"
}

variable "core-config-b64" {
  type      = string
  sensitive = true
}

variable "bot-config-b64" {
  type      = string
  sensitive = true
}

variable "core-env-json-b64" {
  type      = string
  sensitive = true
}

variable "bot-env-json-b64" {
  type      = string
  sensitive = true
}
