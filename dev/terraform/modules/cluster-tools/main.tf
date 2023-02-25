terraform {
  required_providers {
    argocd = {
      source = "oboukili/argocd"
    }
  }
}

resource "kubernetes_namespace" "tools" {
  metadata {
    name = var.namespace
  }
}

resource "argocd_project" "tools" {
  metadata {
    name      = var.argocd_project
    namespace = var.argocd_namespace
    # finalizers = ["resources-finalizer.argocd.argoproj.io"]
  }

  spec {
    source_repos = ["*"]
    destination {
      server    = "https://kubernetes.default.svc"
      namespace = var.namespace
    }
    cluster_resource_whitelist {
      group = "*"
      kind  = "*"
    }
  }
}
