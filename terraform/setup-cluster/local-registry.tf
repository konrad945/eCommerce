resource "kubernetes_config_map" "local-registry" {
  metadata {
    name = "local-registry-hosting"
    namespace = "kube-public"
  }

  data = {
    "localRegistryHosting.v1" = <<EOF
host: "localhost:5001"
help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
    EOF
  }
}