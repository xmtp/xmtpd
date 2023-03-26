locals {
  prometheus_server_endpoint = "prometheus-server:80"
  prometheus_server_url      = "http://${local.prometheus_server_endpoint}"
}

module "argocd_app_prometheus" {
  count  = var.enable_monitoring ? 1 : 0
  source = "../argocd-application"

  argocd_namespace = var.argocd_namespace
  argocd_project   = module.argocd_project.name
  name             = "prometheus"
  namespace        = var.namespace
  wait             = var.wait_for_ready
  repo_url         = "https://prometheus-community.github.io/helm-charts"
  chart            = "prometheus"
  target_revision  = "19.7.2"
  helm_values = [
    <<EOF
      server:
        nodeSelector:
          node-pool: ${var.node_pool}
        persistentVolume:
          enabled: false
        ingress:
          enabled: true
          hosts: ${jsonencode(var.prometheus_hostnames)}
        global:
          evaluation_interval: 30s
          scrape_interval: 10s
          scrape_timeout: 5s
      alertmanager:
        persistence:
          enabled: false
        nodeSelector:
          node-pool: ${var.node_pool}
      kube-state-metrics:
        nodeSelector:
          node-pool: ${var.node_pool}
      prometheus-pushgateway:
        nodeSelector:
          node-pool: ${var.node_pool}
    EOF
  ]
}
