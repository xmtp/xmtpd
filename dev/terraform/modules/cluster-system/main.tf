resource "kubernetes_namespace" "system" {
  metadata {
    name = var.namespace
  }
}

resource "helm_release" "traefik" {
  name       = "traefik"
  namespace  = kubernetes_namespace.system.metadata.0.name
  repository = "https://traefik.github.io/charts"
  version    = "21.1.0"
  chart      = "traefik"

  values = [
    <<-EOF
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
  ]
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
