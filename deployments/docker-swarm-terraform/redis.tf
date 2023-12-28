data "docker_registry_image" "redis" {
  name = "redis:7-alpine"
}

resource "docker_image" "redis" {
  name          = data.docker_registry_image.redis.name
  pull_triggers = [data.docker_registry_image.redis.sha256_digest]
  keep_locally  = true
}

resource "docker_service" "redis" {
  name = "${var.app-name}-redis"

  task_spec {
    container_spec {
      image = docker_image.redis.image_id
    }
    networks_advanced {
      name = docker_network.private.id
    }

    # resources {
    #   limits {
    #     memory_bytes = var.memory-limit
    #   }

    #   reservation {
    #     memory_bytes = 16 * 1024 * 1024
    #   }
    # }
  }
}
