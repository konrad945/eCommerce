resource "helm_release" "vault" {
  name  = "vault"
  repository = "https://helm.releases.hashicorp.com"
  chart = "vault"
  version = "0.20.1"
  wait = true

  namespace = var.backend_svc_namespace
  create_namespace = true
}

resource "null_resource" "init_vault" {
  depends_on = [helm_release.vault]

  provisioner "local-exec" {
    command = "sleep 120 && echo $(kubectl exec vault-0 -n $NAMESPACE -- vault operator init  -key-shares=1 -key-threshold=1 -format=json) > cluster-keys.json"

    environment = {
      NAMESPACE = var.backend_svc_namespace
    }
  }
}

resource "null_resource" "unseal_vault" {
  depends_on = [null_resource.init_vault]

  provisioner "local-exec" {
    command = "kubectl exec vault-0 -n $NAMESPACE -- vault operator unseal $(jq -r \".unseal_keys_b64[]\" cluster-keys.json)"

    environment = {
      NAMESPACE = var.backend_svc_namespace
    }
  }
}

resource "helm_release" "jaeger" {
  name  = "jaeger"
  repository = "https://jaegertracing.github.io/helm-charts"
  chart = "jaeger"

  timeout = 600

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