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

output "argocd_project" {
  value = var.argocd_project
}

output "ingress_public_hostname" {
  value = try(kubernetes_service.traefik.status[0].load_balancer[0].ingress[0].hostname, null)
}
