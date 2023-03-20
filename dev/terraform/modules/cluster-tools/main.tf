terraform {
  required_providers {
    argocd = {
      source = "oboukili/argocd"
      version = "4.3.0"
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

module "chat-app" {
  source = "./chat-app"
  count  = var.enable_chat_app ? 1 : 0

  namespace             = var.namespace
  node_pool_label_key   = var.node_pool_label_key
  node_pool_label_value = var.node_pool
  api_url               = var.public_api_url
  hostnames             = [for hostname in var.hostnames : "chat.${hostname}"]
  ingress_class_name    = var.ingress_class_name
}
