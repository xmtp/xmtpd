resource "kubernetes_namespace" "system" {
  metadata {
    name = var.namespace
  }
}

locals {
  namespace = kubernetes_namespace.system.metadata[0].name
}

resource "helm_release" "argocd" {
  name       = "argocd"
  namespace  = local.namespace
  repository = "https://argoproj.github.io/argo-helm"
  version    = "5.23.5"
  chart      = "argo-cd"

  values = [
    <<-EOF
      server:
        ingress:
          enabled: true
          hosts: ${jsonencode(var.argocd_hostnames)}
          ingressClassName: ${var.ingress_class_name}
      configs:
        cm:
          resource.customizations: |
            networking.k8s.io/Ingress:
              health.lua: |
                hs = {}
                hs.status = "Healthy"
                return hs
        params:
          server.insecure: true
    EOF,
    <<-EOF
    server:
      nodeSelector:
        ${var.node_pool_label_key}: ${var.node_pool}
    dex:
      nodeSelector:
        ${var.node_pool_label_key}: ${var.node_pool}
    redis:
      nodeSelector:
        ${var.node_pool_label_key}: ${var.node_pool}
    applicationSet:
      nodeSelector:
        ${var.node_pool_label_key}: ${var.node_pool}
    notifications:
      nodeSelector:
        ${var.node_pool_label_key}: ${var.node_pool}
    controller:
      nodeSelector:
        ${var.node_pool_label_key}: ${var.node_pool}
    repoServer:
      nodeSelector:
        ${var.node_pool_label_key}: ${var.node_pool}
    EOF
  ]
}

data "kubernetes_secret" "argocd-initial-admin-secret" {
  depends_on = [helm_release.argocd]
  metadata {
    name      = "argocd-initial-admin-secret"
    namespace = local.namespace
  }
}

module "argocd_project" {
  source = "../argocd-project"

  name      = var.argocd_project
  namespace = local.namespace
  destinations = [
    {
      server    = "https://kubernetes.default.svc"
      namespace = local.namespace
    }
  ]
}

module "argocd_app_traefik" {
  source = "../argocd-application"

  argocd_namespace = local.namespace
  argocd_project   = module.argocd_project.name
  name             = "traefik"
  namespace        = local.namespace
  repo_url         = "https://traefik.github.io/charts"
  chart            = "traefik"
  target_revision  = "21.1.0"
  helm_values = concat(
    [
      <<EOF
        service:
          type: ${var.ingress_service_type}
        nodeSelector:
          ${var.node_pool_label_key}: ${var.node_pool}
        providers:
          kubernetesCRD:
            ingressClass: ${var.ingress_class_name}
          kubernetesIngress:
            ingressClass: ${var.ingress_class_name}
            publishedService:
              enabled: true
      EOF
    ],
    var.cluster_http_node_port != null && var.cluster_https_node_port != null ? [
      <<EOF
        ports:
          web:
            nodePort: ${var.cluster_http_node_port}
          websecure:
            nodePort: ${var.cluster_https_node_port}
      EOF
    ] : [],
  )
}

data "kubernetes_service" "traefik_service" {
  depends_on = [module.argocd_app_traefik]
  metadata {
    name      = "traefik"
    namespace = local.namespace
  }
}
