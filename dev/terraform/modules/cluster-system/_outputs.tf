output "namespace" {
  value = kubernetes_namespace.system.metadata[0].name
}

output "argocd_hostnames" {
  value = var.argocd_hostnames
}

output "argocd_username" {
  value = "admin"
}

output "argocd_password" {
  value     = data.kubernetes_secret.argocd-initial-admin-secret.data.password
  sensitive = true
}

output "ingress_public_hostname" {
  value = try(one(one(one(data.kubernetes_service.traefik_service.status[*]).load_balancer[*]).ingress[*]).hostname, null)
}
