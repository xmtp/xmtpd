terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
    helm = {
      source = "hashicorp/helm"
    }
  }
}

data "digitalocean_kubernetes_versions" "cluster_versions" {
  version_prefix = "1."
}

data "digitalocean_kubernetes_cluster" "cluster" {
  name       = var.cluster_name
  depends_on = [digitalocean_kubernetes_cluster.cluster]
}

resource "digitalocean_kubernetes_cluster" "cluster" {
  name         = var.cluster_name
  region       = var.cluster_region
  auto_upgrade = true
  version      = data.digitalocean_kubernetes_versions.cluster_versions.latest_version

  node_pool {
    name       = var.xmtp_utils_pool_label_value
    size       = var.xmtp_utils_pool_node_size
    node_count = var.num_xmtp_utils_pool_nodes
    labels = {
      (var.node_pool_label_key) = var.xmtp_utils_pool_label_value
    }
  }
}

resource "digitalocean_kubernetes_node_pool" "xmtp-nodes" {
  cluster_id = digitalocean_kubernetes_cluster.cluster.id

  name       = var.xmtp_nodes_pool_label_value
  size       = var.xmtp_nodes_pool_node_size
  node_count = var.num_xmtp_nodes_pool_nodes
  labels = {
    (var.node_pool_label_key) = var.xmtp_nodes_pool_label_value
  }
}
