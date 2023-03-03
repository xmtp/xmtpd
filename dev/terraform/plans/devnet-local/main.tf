terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
    helm = {
      source = "hashicorp/helm"
    }
    argocd = {
      source = "oboukili/argocd"
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

provider "argocd" {
  server_addr                 = module.system.argocd_hostnames[0]
  username                    = module.system.argocd_username
  password                    = module.system.argocd_password
  plain_text                  = true
  insecure                    = true
  port_forward                = true
  port_forward_with_namespace = module.system.namespace
}

locals {
  node_pool_label_key     = "node-pool"
  system_node_pool        = "xmtp-system"
  nodes_node_pool         = "xmtp-nodes"
  ingress_class_name      = "traefik"
  cluster_http_node_port  = 32080
  cluster_https_node_port = 32443
  hostnames               = ["localhost", "xmtp.local"]

  nodes = [for node in var.nodes : {
    name = node.name
  }]
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
}

module "nodes" {
  source     = "../../modules/cluster-nodes"
  depends_on = [module.system]

  namespace           = "xmtp-nodes"
  node_pool_label_key = local.node_pool_label_key
  node_pool           = local.nodes_node_pool
  argocd_namespace    = module.system.namespace
  nodes               = local.nodes
  argocd_project      = "xmtp-nodes"
  ingress_class_name  = local.ingress_class_name
}
