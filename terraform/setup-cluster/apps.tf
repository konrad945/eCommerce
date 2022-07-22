resource "helm_release" "catalog" {
  name  = "catalog"
  chart = "../../svc/catalog/helm"

  namespace = var.backend_svc_namespace
  create_namespace = true

  depends_on = [helm_release.jaeger, helm_release.postgresql, null_resource.vault_role]
}