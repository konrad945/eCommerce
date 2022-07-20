provider "kubernetes" {
  config_path = pathexpand(var.kind_cluster_config_path)
  config_context_cluster = "kind-cluster"
}

provider "helm" {
  kubernetes {
    config_path = pathexpand(var.kind_cluster_config_path)
  }
}