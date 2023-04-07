output "k8s_cluster_name" {
  value = module.cluster.k8s_cluster_name
}

output "nodes" {
  value = var.nodes
}
