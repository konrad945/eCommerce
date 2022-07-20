resource "helm_release" "postgresql" {
  name = "postgre-db"
  repository = "https://charts.bitnami.com/bitnami"
  chart = "postgresql"

  namespace = var.backend_svc_namespace
  create_namespace = true

  values = [file(pathexpand("../../helm/postgresql/values.yaml"))]
}