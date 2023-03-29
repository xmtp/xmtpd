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
    },
    {
      server    = "https://kubernetes.default.svc"
      namespace = "kube-system"
    },
  ]
}

resource "kubernetes_service" "traefik" {
  depends_on             = [module.argocd_app_traefik]
  wait_for_load_balancer = var.ingress_service_type == "LoadBalancer"
  metadata {
    name      = "traefik"
    namespace = local.namespace
    labels = {
      "app.kubernetes.io/name" = "traefik"
    }
    annotations = {}
  }
  spec {
    type = var.ingress_service_type
    selector = {
      "app.kubernetes.io/instance" = "traefik-${local.namespace}"
      "app.kubernetes.io/name"     = "traefik"
    }
    port {
      name        = "web"
      port        = 80
      protocol    = "TCP"
      target_port = "web"
      node_port   = var.cluster_http_node_port
    }
    port {
      name        = "websecure"
      port        = 443
      protocol    = "TCP"
      target_port = "websecure"
      node_port   = var.cluster_https_node_port
    }
  }
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
  helm_values = [
    <<EOF
      service:
        enabled: false
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
  ]
}
