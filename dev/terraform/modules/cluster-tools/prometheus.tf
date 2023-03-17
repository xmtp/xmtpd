locals {
  prometheus_hostnames       = [for hostname in var.hostnames : "prometheus.${hostname}"]
  prometheus_public_hostname = local.prometheus_hostnames[0]
  prometheus_server_endpoint = "prometheus-server:80"
  prometheus_server_url      = "http://${local.prometheus_server_endpoint}"
}

resource "argocd_application" "prometheus" {
  count      = var.enable_monitoring ? 1 : 0
  depends_on = [argocd_project.tools]
  wait       = var.wait_for_ready
  metadata {
    name      = "prometheus"
    namespace = var.argocd_namespace
  }
  spec {
    project = argocd_project.tools.metadata[0].name
    source {
      repo_url        = "https://prometheus-community.github.io/helm-charts"
      chart           = "prometheus"
      target_revision = "19.7.2"
      helm {
        release_name = "prometheus"
        values       = <<EOT
          server:
            nodeSelector:
              node-pool: ${var.node_pool}
            persistentVolume:
              enabled: false
            ingress:
              enabled: true
              hosts:
                - ${local.prometheus_public_hostname}
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
        EOT
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
