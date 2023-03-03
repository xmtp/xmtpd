resource "argocd_application" "metrics-server" {
  metadata {
    name      = "metrics-server"
    namespace = var.argocd_namespace
    # finalizers = ["resources-finalizer.argocd.argoproj.io"]
  }
  spec {
    project = argocd_project.tools.metadata.0.name
    source {
      repo_url        = "https://kubernetes-sigs.github.io/metrics-server/"
      chart           = "metrics-server"
      target_revision = "3.8.3"
      helm {
        release_name = "metrics-server"
        values       = <<EOT
args:
  - --kubelet-insecure-tls
  - --kubelet-preferred-address-types=InternalIP
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
        allow_empty = false
      }
    }
  }
}
