locals {
  num_nodes       = length(var.nodes)
  nodes_mid_index = floor(local.num_nodes / 2)
  nodes_group1    = slice(var.nodes, 0, local.nodes_mid_index)
  nodes_group2    = slice(var.nodes, local.nodes_mid_index, local.num_nodes)
  namespace       = kubernetes_namespace.nodes.metadata[0].name
}

resource "kubernetes_namespace" "nodes" {
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
      namespace = local.namespace
    }
  ]
}

module "nodes_group1" {
  source     = "./node"
  depends_on = [kubernetes_namespace.nodes]
  count      = length(local.nodes_group1)

  name                      = local.nodes_group1[count.index].name
  namespace                 = local.namespace
  argocd_project            = var.argocd_project
  argocd_namespace          = var.argocd_namespace
  p2p_persistent_peers      = local.nodes_group1[count.index].p2p_persistent_peers
  private_key               = var.node_keys[local.nodes_group1[count.index].name]
  container_image           = var.container_image
  storage_class_name        = var.storage_class_name
  storage_request           = var.container_storage_request
  cpu_request               = var.container_cpu_request
  hostnames                 = [for hostname in var.hostnames : "${local.nodes_group1[count.index].name}.${hostname}"]
  p2p_port                  = 9000
  api_grpc_port             = 5000
  api_http_port             = 5001
  metrics_port              = 8009
  node_pool_label_key       = var.node_pool_label_key
  node_pool                 = var.node_pool
  one_instance_per_k8s_node = var.one_instance_per_k8s_node
  ingress_class_name        = var.ingress_class_name
  wait_for_ready            = var.wait_for_ready
  debug                     = var.debug
  store_type                = local.nodes_group1[count.index].store_type
}

module "nodes_group2" {
  source     = "./node"
  depends_on = [kubernetes_namespace.nodes, module.nodes_group1]
  count      = length(local.nodes_group2)

  name                      = local.nodes_group2[count.index].name
  namespace                 = local.namespace
  argocd_project            = var.argocd_project
  argocd_namespace          = var.argocd_namespace
  p2p_persistent_peers      = local.nodes_group2[count.index].p2p_persistent_peers
  private_key               = var.node_keys[local.nodes_group2[count.index].name]
  container_image           = var.container_image
  storage_class_name        = var.storage_class_name
  storage_request           = var.container_storage_request
  cpu_request               = var.container_cpu_request
  hostnames                 = [for hostname in var.hostnames : "${local.nodes_group2[count.index].name}.${hostname}"]
  p2p_port                  = 9000
  api_grpc_port             = 5000
  api_http_port             = 5001
  metrics_port              = 8009
  node_pool_label_key       = var.node_pool_label_key
  node_pool                 = var.node_pool
  one_instance_per_k8s_node = var.one_instance_per_k8s_node
  ingress_class_name        = var.ingress_class_name
  wait_for_ready            = var.wait_for_ready
  debug                     = var.debug
  store_type                = local.nodes_group2[count.index].store_type
}

resource "kubernetes_service" "nodes_api" {
  metadata {
    name      = "nodes-api"
    namespace = local.namespace
  }
  spec {
    selector = {
      "app.kubernetes.io/part-of" = "xmtp-nodes"
    }
    port {
      name        = "http"
      port        = var.node_api_http_port
      target_port = var.node_api_http_port
    }
  }
}

resource "kubernetes_ingress_v1" "nodes_api" {
  metadata {
    name      = "nodes-api"
    namespace = local.namespace
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
                name = kubernetes_service.nodes_api.metadata[0].name
                port {
                  number = kubernetes_service.nodes_api.spec[0].port[0].port
                }
              }
            }
          }
        }
      }
    }
  }
}
