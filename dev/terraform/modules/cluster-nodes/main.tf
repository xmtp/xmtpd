terraform {
  required_providers {
    argocd = {
      source = "oboukili/argocd"
    }
  }
}

resource "kubernetes_namespace" "nodes" {
  metadata {
    name = var.namespace
  }
}

resource "argocd_project" "nodes" {
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

resource "argocd_application" "nodes" {
  count = length(var.nodes)
  metadata {
    name      = var.nodes[count.index].name
    namespace = var.argocd_namespace
    # finalizers = ["resources-finalizer.argocd.argoproj.io"]
  }
  spec {
    project = argocd_project.nodes.metadata.0.name
    source {
      repo_url = "https://github.com/argoproj/argocd-example-apps.git"
      path     = "helm-guestbook"
      # repo_url = "https://github.com/xmtp/xmtpd.git"
      # path = "dev/helm/xmtp-node"
      target_revision = "HEAD"
      helm {
        values = yamlencode({
          ingress = {
            enabled = true
            hosts = [
              {
                // TODO: use var hostnames list to build this
                host = "${var.nodes[count.index].name}.localhost"
                paths = [
                  {
                    path     = "/"
                    pathType = "Prefix"
                  }
                ]
              }
            ]
          }
          nodeSelector = {
            // TODO: use label key var
            "node-pool" = var.node_pool
          }
        })
      }
    }

    destination {
      server    = "https://kubernetes.default.svc"
      namespace = var.namespace
    }

    // TODO: disable this just for local dev clusters
    # sync_policy {
    #   automated = {
    #     prune     = true
    #     self_heal = true
    #     allow_empty = false
    #   }
    # }
  }
}
