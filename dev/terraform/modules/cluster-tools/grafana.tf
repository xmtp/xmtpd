resource "argocd_application" "grafana" {
  metadata {
    name      = "grafana"
    namespace = var.argocd_namespace
    # finalizers = ["resources-finalizer.argocd.argoproj.io"]
  }
  spec {
    project = argocd_project.tools.metadata.0.name
    source {
      repo_url        = "https://grafana.github.io/helm-charts"
      chart           = "grafana"
      target_revision = "6.51.2"
      helm {
        release_name = "grafana"
        values       = <<EOT
nodeSelector:
  node-pool: ${var.node_pool}
persistence:
  enabled: false
ingress:
  enabled: true
  hosts:
    - grafana.localhost
grafana.ini:
  auth.anonymous:
    enabled: true
    org_name: "Main Org."
    # Role for unauthenticated users, other valid values are `Editor` and `Admin`
    org_role: "Editor"
    hide_version: true
EOT
      }
    }

    destination {
      server    = "https://kubernetes.default.svc"
      namespace = var.namespace
    }

    sync_policy {
      automated = {
        prune     = true
        self_heal = true
      }
    }
  }
}
