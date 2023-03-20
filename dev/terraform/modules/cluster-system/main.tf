terraform {
  required_providers {
    argocd = {
      source = "oboukili/argocd"
      version = "4.3.0"
    }
  }
}

resource "kubernetes_namespace" "system" {
  metadata {
    name = var.namespace
  }
}

resource "helm_release" "argocd" {
  name       = "argocd"
  namespace  = var.namespace
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
    namespace = var.namespace
  }
}

resource "argocd_project" "system" {
  metadata {
    name      = var.argocd_project
    namespace = var.namespace
  }

  spec {
    source_repos = ["*"]
    destination {
      server    = "https://kubernetes.default.svc"
      namespace = var.namespace
    }
    cluster_resource_whitelist {
      group = "*"
      kind  = "*"
    }
  }
}

resource "argocd_application" "traefik" {
  depends_on = [helm_release.argocd]
  wait       = true
  metadata {
    name      = "traefik"
    namespace = var.namespace
  }
  spec {
    project = argocd_project.system.metadata[0].name
    source {
      repo_url        = "https://traefik.github.io/charts"
      chart           = "traefik"
      target_revision = "21.1.0"
      helm {
        release_name = "traefik"
        values       = <<EOF
          service:
            type: NodePort
          nodeSelector:
            ${var.node_pool_label_key}: ${var.node_pool}
          ports:
            web:
              nodePort: ${var.cluster_http_node_port}
            websecure:
              nodePort: ${var.cluster_https_node_port}
          providers:
            kubernetesCRD:
              ingressClass: ${var.ingress_class_name}
            kubernetesIngress:
              ingressClass: ${var.ingress_class_name}
              publishedService:
                enabled: true
        EOF
      }
    }

    destination {
      server    = "https://kubernetes.default.svc"
      namespace = var.namespace
    }

    sync_policy {
      automated = {
        prune       = true
        self_heal   = true
        allow_empty = false
      }
    }
  }
}
