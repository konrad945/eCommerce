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
    command = "sleep 80 && echo $(kubectl exec vault-0 -n $NAMESPACE -- vault operator init  -key-shares=1 -key-threshold=1 -format=json) > cluster-keys.json"

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

resource "null_resource" "vault_login" {
  depends_on = [null_resource.unseal_vault]

  provisioner "local-exec" {
    command = "kubectl exec vault-0 -n $NAMESPACE -- vault login $(jq -r \".root_token\" cluster-keys.json)"

    environment = {
      NAMESPACE = var.backend_svc_namespace
    }
  }
}

resource "null_resource" "vault_database-enable" {
  depends_on = [null_resource.vault_login]

  provisioner "local-exec" {
    command = <<-TOML
      kubectl exec vault-0 -n $NAMESPACE -- vault secrets enable database
    TOML

    environment = {
      NAMESPACE = var.backend_svc_namespace
    }
  }
}

resource "null_resource" "vault_database-roles" {
  depends_on = [null_resource.vault_database-enable]

  provisioner "local-exec" {
    command = <<-TOML
      kubectl exec vault-0 -n $NAMESPACE -- vault write database/roles/db-app \
        db_name=catalog \
        creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; \
          GRANT SELECT ON ALL TABLES IN SCHEMA public TO \"{{name}}\";" \
        revocation_statements="ALTER ROLE \"{{name}}\" NOLOGIN;"\
        default_ttl="1h" \
        max_ttl="24h"
    TOML

    environment = {
      NAMESPACE = var.backend_svc_namespace
    }
  }
}

resource "null_resource" "vault_database_plugin_postgre" {
  depends_on = [null_resource.vault_database-roles]

  provisioner "local-exec" {
    command = <<-TOML
      kubectl exec vault-0 -n $NAMESPACE -- vault write database/config/catalog \
        plugin_name=postgresql-database-plugin \
        allowed_roles="*" \
        connection_url="postgresql://{{username}}:{{password}}@postgre-db-postgresql:5432/catalog?sslmode=disable" \
        username="postgres" \
        password="password"
    TOML

    environment = {
      NAMESPACE = var.backend_svc_namespace
    }
  }
}

resource "null_resource" "vault_database_catalog-rotate_password" {
  depends_on = [null_resource.vault_database_plugin_postgre]

  provisioner "local-exec" {
    command = <<-TOML
      kubectl exec vault-0 -n $NAMESPACE -- vault write --force /database/rotate-root/catalog
    TOML

    environment = {
      NAMESPACE = var.backend_svc_namespace
    }
  }
}

resource "null_resource" "vault_kubernetes_enable" {
  depends_on = [null_resource.vault_database_plugin_postgre]

  provisioner "local-exec" {
    command = <<-TOML
      kubectl exec vault-0 -n $NAMESPACE -- vault auth enable kubernetes
    TOML

    environment = {
      NAMESPACE = var.backend_svc_namespace
    }
  }
}

resource "null_resource" "vault_kubernetes-configure1" {
  depends_on = [null_resource.vault_kubernetes_enable]

  provisioner "local-exec" {
    command = <<-TOML
      kubectl exec vault-0 -n backend -- sh -c 'vault write auth/kubernetes/config \
       token_reviewer_jwt="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
       kubernetes_host=https://$${KUBERNETES_PORT_443_TCP_ADDR}:443 \
       kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt'
    TOML

    environment = {
      NAMESPACE = var.backend_svc_namespace
    }
  }
}

resource "null_resource" "vault_policy" {
  depends_on = [null_resource.vault_kubernetes-configure1]

  provisioner "local-exec" {
    command = <<-TOML
      kubectl cp ../../vault $NAMESPACE/vault-0:/tmp
      kubectl exec vault-0 -n $NAMESPACE -- vault policy write catalog /tmp/vault/catalog/policy.hcl
    TOML

    environment = {
      NAMESPACE = var.backend_svc_namespace
    }
  }
}

resource "null_resource" "vault_role" {
  depends_on = [null_resource.vault_policy]

  provisioner "local-exec" {
    command = <<-TOML
      kubectl exec vault-0 -n $NAMESPACE -- vault write auth/kubernetes/role/catalog \
        bound_service_account_names=catalog \
        bound_service_account_namespaces=backend \
        policies=catalog \
        ttl=10h
    TOML

    environment = {
      NAMESPACE = var.backend_svc_namespace
    }
  }
}