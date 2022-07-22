resource "helm_release" "jaeger" {
  name  = "jaeger"
  repository = "https://jaegertracing.github.io/helm-charts"
  chart = "jaeger"

  namespace = var.backend_svc_namespace
  create_namespace = true

  values = [file(pathexpand("../../helm/jaeger/values.yaml"))]
}

resource "helm_release" "postgresql" {
  name = "postgre-db"
  repository = "https://charts.bitnami.com/bitnami"
  chart = "postgresql"

  namespace = var.backend_svc_namespace
  create_namespace = true

  values = [file(pathexpand("../../helm/postgresql/values.yaml"))]
}

resource "helm_release" "elastic" {
  name = "efk-elastic"
  repository = "https://helm.elastic.co"
  chart = "elasticsearch"
  version = "7.17.3"

  namespace = var.backend_svc_namespace
  create_namespace = true

  values = [file(pathexpand("../../helm/elastic/values.yaml"))]
}

resource "helm_release" "fluentd" {
  name = "efk-fluentd"
  repository = "https://fluent.github.io/helm-charts"
  chart = "fluentd"
  version = "0.3.9"

  namespace = var.backend_svc_namespace
  create_namespace = true

  values = [file(pathexpand("../../helm/fluentd/values.yaml"))]

  depends_on = [helm_release.elastic]
}

resource "helm_release" "kibana" {
  name = "efk-kibana"
  repository = "https://helm.elastic.co"
  chart = "kibana"
  version = "7.17.3"

  namespace = var.backend_svc_namespace
  create_namespace = true

  values = [file(pathexpand("../../helm/kibana/values.yaml"))]

  depends_on = [helm_release.elastic]
}