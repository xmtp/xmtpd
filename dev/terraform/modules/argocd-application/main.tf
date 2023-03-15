resource "kubernetes_manifest" "argocd_application" {
  computed_fields = [
    "metadata.labels",
    "metadata.annotations",
    "metadata.finalizers",
    "spec.source.helm.version"
  ]
  field_manager {
    force_conflicts = true
  }
  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "Application"
    metadata = {
      name      = var.name
      namespace = var.argocd_namespace
    }
    spec = merge(
      {
        project = var.argocd_project
        source = {
          repoURL        = var.repo_url
          targetRevision = var.target_revision
          chart          = var.chart
          path           = var.path
          helm = {
            releaseName = var.name
            values      = yamlencode(merge([for values in var.helm_values : yamldecode(values)]...))
          }
        }
        destination = {
          server    = var.destination_server
          namespace = var.namespace
        }
      },
      var.auto_sync ? {
        syncPolicy = {
          automated = {
            prune    = true
            selfHeal = true
          }
        }
      } : {}
    )
  }
  dynamic "wait" {
    for_each = var.wait ? [1] : []
    content {
      fields = {
        "status.health.status" = "Healthy"
      }
    }
  }
}
