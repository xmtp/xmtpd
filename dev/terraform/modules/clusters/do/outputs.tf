output "cluster_endpoint" {
  value = data.digitalocean_kubernetes_cluster.cluster.endpoint
}

output "cluster_ca_certificate" {
  value = base64decode(
    data.digitalocean_kubernetes_cluster.cluster.kube_config[0].cluster_ca_certificate
  )
}

output "cluster_id" {
  value = data.digitalocean_kubernetes_cluster.cluster.id
}
