resource "kind_cluster" "default" {
  name = var.kind_custer_name
  kubeconfig_path = pathexpand(var.kind_cluster_config_path)
  wait_for_ready = true

  depends_on = [docker_container.local_registry]

  kind_config {
    kind        = "cluster"
    api_version = "kind.x-k8s.io/v1alpha4"
    containerd_config_patches = [
      <<-TOML
      [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:5001"]
          endpoint = ["http://kind-registry:5000"]
      TOML
    ]
  }

  provisioner "local-exec" {
    command = <<EOF
    #!/bin/sh
    set -o errexit
    docker network connect "kind" "kind-registry"
    EOF
  }
}