locals {
  labels = {
    "app.kubernetes.io/name" = var.name
  }
}

resource "kubernetes_service" "service" {
  metadata {
    name      = var.name
    namespace = var.namespace
  }
  spec {
    selector = local.labels
    port {
      name        = "web"
      port        = var.service_port
      target_port = var.container_port
    }
  }
}

resource "kubernetes_deployment" "deployment" {
  metadata {
    name      = var.name
    namespace = var.namespace
    labels    = local.labels
  }
  spec {
    replicas = 1
    selector {
      match_labels = local.labels
    }
    template {
      metadata {
        labels = local.labels
      }
      spec {
        node_selector = {
          (var.node_pool_label_key) = var.node_pool_label_value
        }
        container {
          name  = "web"
          image = var.container_image
          port {
            name           = "http"
            container_port = var.container_port
          }
          env {
            name  = "PORT"
            value = var.container_port
          }
          env {
            name  = "NEXT_PUBLIC_XMTP_API_URL"
            value = var.api_url
          }
        }
      }
    }
  }
}

resource "kubernetes_ingress_v1" "app" {
  count = length(var.hostnames) > 0 ? 1 : 0
  metadata {
    name      = var.name
    namespace = var.namespace
  }
  spec {
    ingress_class_name = var.ingress_class_name
    dynamic "rule" {
      for_each = var.hostnames
      content {
        host = rule.value
        http {
          path {
            path = "/"
            backend {
              service {
                name = var.name
                port {
                  number = var.service_port
                }
              }
            }
          }
        }
      }
    }
  }
}
