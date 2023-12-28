resource "docker_network" "private" {
  name   = "${var.app-name}-network"
  driver = "overlay"
}
