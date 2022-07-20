resource "docker_container" "local_registry" {
  name = "kind-registry"
  image = "registry:2"

  ports {
    ip = "127.0.0.1"
    internal = 5000
    external = 5001
  }
}