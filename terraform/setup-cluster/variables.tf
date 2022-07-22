variable "kind_cluster_config_path" {
  type = string
  description = "The location where this cluster's kubeconfig will be saved to."
  default = "~/.kube/config"
}

variable "backend_svc_namespace" {
  type = string
  description = "The backend services namespace (it will be created)"
  default = "backend"
}

variable "apps_svc_namespace" {
  type = string
  description = "The apps services namespace (it will be created)"
  default = "apps"
}