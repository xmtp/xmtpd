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
    spec = {
      description                = var.description
      sourceRepos                = var.source_repos
      destinations               = var.destinations
      clusterResourceWhitelist   = var.cluster_resource_whitelist
      namespaceResourceWhitelist = var.namespace_resource_whitelist
      namespaceResourceBlacklist = var.namespace_resource_blacklist
    }
  }
}

resource "time_sleep" "wait_after_create" {
  depends_on = [kubernetes_manifest.argocd_project]

  create_duration = "2s"
}
