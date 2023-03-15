resource "kubernetes_manifest" "argocd_project" {
  computed_fields = [
    "metadata.labels",
    "metadata.annotations",
    "metadata.finalizers"
  ]
  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "AppProject"
    metadata = {
      name      = var.name
      namespace = var.namespace
    }
    spec = merge(
      {
        description  = var.description
        sourceRepos  = var.source_repos
        destinations = var.destinations
      },
      var.cluster_resource_whitelist == null ? {} : { clusterResourceWhitelist = var.cluster_resource_whitelist },
      var.namespace_resource_whitelist == null ? {} : { namespaceResourceWhitelist = var.namespace_resource_whitelist },
      var.namespace_resource_blacklist == null ? {} : { namespaceResourceBlacklist = var.namespace_resource_blacklist },
    )
  }
}

resource "time_sleep" "wait_after_create" {
  depends_on = [kubernetes_manifest.argocd_project]

  create_duration = "2s"
}
