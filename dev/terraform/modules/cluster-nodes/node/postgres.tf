locals {
  postgres_name         = "${var.name}-postgres"
  postgres_password     = one(data.kubernetes_secret.postgres[*].data.postgres-password)
  postgres_service_name = "${local.postgres_name}-postgresql"
  postgres_dsn          = local.postgres_password != null ? "postgres://postgres:${local.postgres_password}@${local.postgres_service_name}:5432?sslmode=disable" : null
}

resource "argocd_application" "postgres" {
  count = var.enable_postgres ? 1 : 0
  wait  = true
  metadata {
    name      = local.postgres_name
    namespace = var.argocd_namespace
  }
  spec {
    project = var.argocd_project
    source {
      repo_url        = "https://charts.bitnami.com/bitnami"
      chart           = "postgresql"
      target_revision = "12.2.3"
      helm {
        release_name = local.postgres_name
        values       = <<EOT
        EOT
      }
    }

    destination {
      server    = "https://kubernetes.default.svc"
      namespace = var.namespace
    }

    sync_policy {
      automated = {
        prune       = true
        self_heal   = true
        allow_empty = false
      }
    }
  }
}

data "kubernetes_secret" "postgres" {
  count      = var.enable_postgres ? 1 : 0
  depends_on = [argocd_application.postgres]
  metadata {
    name      = "${local.postgres_name}-postgresql"
    namespace = var.namespace
  }
}
