locals {
  postgres_name         = "${var.name}-postgres"
  postgres_password     = one(data.kubernetes_secret.postgres[*].data.postgres-password)
  postgres_service_name = "${local.postgres_name}-postgresql"
  postgres_dsn          = local.postgres_password != null ? "postgres://postgres:${local.postgres_password}@${local.postgres_service_name}:5432?sslmode=disable" : null
}

module "argocd_app_postgres" {
  count  = var.store_type == "postgres" ? 1 : 0
  source = "../../argocd-application"

  argocd_namespace = var.argocd_namespace
  argocd_project   = var.argocd_project
  name             = local.postgres_name
  namespace        = var.namespace
  repo_url         = "https://charts.bitnami.com/bitnami"
  chart            = "postgresql"
  target_revision  = "12.2.3"
  wait             = true
}

data "kubernetes_secret" "postgres" {
  count      = var.store_type == "postgres" ? 1 : 0
  depends_on = [module.argocd_app_postgres]
  metadata {
    name      = "${local.postgres_name}-postgresql"
    namespace = var.namespace
  }
}
