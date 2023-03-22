terraform {
  required_providers {
    kind = {
      source  = "tehcyx/kind"
      version = "0.0.16"
    }
  }
}

resource "kind_cluster" "cluster" {
  name            = var.name
  kubeconfig_path = var.kubeconfig_path
  wait_for_ready  = var.wait_for_ready
  kind_config {
    kind        = "Cluster"
    api_version = "kind.x-k8s.io/v1alpha4"
    node {
      role = "control-plane"
    }
    dynamic "node" {
      for_each = var.nodes
      content {
        role   = "worker"
        labels = coalesce(node.value.labels, {})
        dynamic "extra_port_mappings" {
          for_each = [for k, v in coalesce(node.value.extra_port_mappings, {}) : { container_port = k, host_port = v }]
          content {
            container_port = extra_port_mappings.value.container_port
            host_port      = extra_port_mappings.value.host_port
          }
        }
      }
    }
  }
}
