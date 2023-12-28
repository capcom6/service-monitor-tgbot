variable "service-name" {
  type = string
}
variable "service-version" {
  type = string
}

variable "config-b64" {
  type = string
}
variable "env-json-b64" {
  type = string
}

variable "network-id" {
  type = string
}

variable "memory-limit" {
  default = 32 * 1024 * 1024
}

variable "memory-reserve" {
  default = 16 * 1024 * 1024
}

data "docker_registry_image" "app-image" {
  name = "capcom6/${var.service-name}:${var.service-version}"
}

resource "docker_image" "app" {
  name          = data.docker_registry_image.app-image.name
  pull_triggers = [data.docker_registry_image.app-image.sha256_digest]
  keep_locally  = true
}

resource "docker_config" "app" {
  name = "${var.service-name}-config.yml-${replace(timestamp(), ":", ".")}"
  data = var.config-b64

  lifecycle {
    ignore_changes        = [name]
    create_before_destroy = true
  }
}

resource "docker_service" "app" {
  name = var.service-name

  task_spec {
    container_spec {
      image = docker_image.app.image_id

      configs {
        config_id   = docker_config.app.id
        config_name = docker_config.app.name
        file_name   = "/app/config.yml"
        file_uid    = 405
        file_gid    = 100
      }

      env = jsondecode(base64decode(var.env-json-b64))
    }
    networks_advanced {
      name = var.network-id
    }

    resources {
      limits {
        memory_bytes = var.memory-limit
      }

      reservation {
        memory_bytes = 16 * 1024 * 1024
      }
    }
  }
}
