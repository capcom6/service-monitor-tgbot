data "docker_registry_image" "app-image" {
  name = "capcom6/${var.app-name}:${var.app-version}"
}

data "docker_network" "proxy" {
  name = "proxy"
}


resource "docker_image" "app" {
  name          = data.docker_registry_image.app-image.name
  pull_triggers = [data.docker_registry_image.app-image.sha256_digest]
  keep_locally  = true
}

resource "docker_config" "app" {
  name = "${var.app-name}-config.yml-${replace(timestamp(), ":", ".")}"
  data = var.app-config-b64

  lifecycle {
    ignore_changes        = [name]
    create_before_destroy = true
  }
}

resource "docker_service" "app" {
  name = var.app-name

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

      env = jsondecode(base64decode(var.app-env-json-b64))
    }
    networks_advanced {
      name = data.docker_network.proxy.id
    }

    resources {
      limits {
        # nano_cpus    = var.cpu-limit
        memory_bytes = var.memory-limit
      }

      reservation {
        # nano_cpus    = 10 * 10000000
        memory_bytes = 16 * 1024 * 1024
      }
    }
  }
}
