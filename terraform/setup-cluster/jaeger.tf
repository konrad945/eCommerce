resource "helm_release" "jaeger" {
  name  = "jaeger"
  repository = "https://jaegertracing.github.io/helm-charts"
  chart = "jaeger"

  namespace = var.backend_svc_namespace
  create_namespace = true

  values = [file(pathexpand("../../helm/jaeger/values.yaml"))]
}