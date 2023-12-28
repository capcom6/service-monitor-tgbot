module "core" {
  source = "./modules/service"

  service-name    = "service-monitor-core"
  service-version = var.app-version

  config-b64   = var.core-config-b64
  env-json-b64 = var.core-env-json-b64

  network-id = docker_network.private.id

  depends_on = [docker_service.redis]
}

module "bot" {
  source = "./modules/service"

  service-name    = "service-monitor-bot"
  service-version = var.app-version

  config-b64   = var.bot-config-b64
  env-json-b64 = var.bot-env-json-b64

  network-id = docker_network.private.id

  depends_on = [docker_service.redis]
}
