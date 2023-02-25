output "namespace" {
  value = kubernetes_namespace.system.metadata.0.name
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
