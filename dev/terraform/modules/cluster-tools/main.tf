resource "kubernetes_namespace" "tools" {
  metadata {
    name = var.namespace
  }
}

module "argocd_project" {
  source = "../argocd-project"

  name      = var.argocd_project
  namespace = var.argocd_namespace
  destinations = [
    {
      server    = "https://kubernetes.default.svc"
      namespace = var.namespace
    }
  ]
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
