terraform {
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "2.9.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.18.1"
    }
  }
}

provider "kubernetes" {
  config_path = var.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig_path
  }
}

locals {
  node_pool_label_key     = "node-pool"
  system_node_pool        = "xmtp-system"
  nodes_node_pool         = "xmtp-nodes"
  ingress_class_name      = "traefik"
  cluster_http_node_port  = 32080
  cluster_https_node_port = 32443
  hostnames               = ["localhost", "xmtp.local"]
  node_api_http_port      = 5001
}

module "cluster" {
  source = "../../modules/clusters/kind"

  name            = "xmtp-devnet-local"
  kubeconfig_path = startswith(var.kubeconfig_path, "/") ? var.kubeconfig_path : abspath(var.kubeconfig_path)
  nodes = concat(
    [{
      labels = {
        (local.node_pool_label_key) = local.system_node_pool
        "ingress-ready"             = "true"
      }
      extra_port_mappings = {
        (local.cluster_http_node_port)  = 80
        (local.cluster_https_node_port) = 443
      }
    }],
    [for i in range(var.num_xmtp_node_pool_nodes) : {
      labels = {
        (local.node_pool_label_key) = local.nodes_node_pool
      }
    }]
  )
}

module "system" {
  source     = "../../modules/cluster-system"
  depends_on = [module.cluster]

  namespace               = "xmtp-system"
  node_pool_label_key     = local.node_pool_label_key
  node_pool               = local.system_node_pool
  argocd_project          = "xmtp-system"
  cluster_http_node_port  = local.cluster_http_node_port
  cluster_https_node_port = local.cluster_https_node_port
  argocd_hostnames        = [for hostname in local.hostnames : "argo.${hostname}"]
  ingress_class_name      = local.ingress_class_name
}

module "tools" {
  source     = "../../modules/cluster-tools"
  depends_on = [module.system]

  namespace           = "xmtp-tools"
  node_pool_label_key = local.node_pool_label_key
  node_pool           = local.system_node_pool
  argocd_namespace    = module.system.namespace
  argocd_project      = "xmtp-tools"
  ingress_class_name  = local.ingress_class_name
  wait_for_ready      = false
  enable_chat_app     = var.enable_chat_app
  enable_monitoring   = var.enable_monitoring
  hostnames           = local.hostnames
  public_api_url      = "http://${local.hostnames[0]}"
}

module "nodes" {
  source     = "../../modules/cluster-nodes"
  depends_on = [module.system]

  namespace                 = "xmtp-nodes"
  container_image           = var.node_container_image
  node_pool_label_key       = local.node_pool_label_key
  node_pool                 = local.nodes_node_pool
  argocd_namespace          = module.system.namespace
  argocd_project            = "xmtp-nodes"
  nodes                     = var.nodes
  node_keys                 = var.node_keys
  ingress_class_name        = local.ingress_class_name
  hostnames                 = local.hostnames
  node_api_http_port        = local.node_api_http_port
  storage_class_name        = "standard"
  container_storage_request = "1Gi"
  container_cpu_request     = "10m"
  debug                     = true
  wait_for_ready            = false
  one_instance_per_k8s_node = false
}
