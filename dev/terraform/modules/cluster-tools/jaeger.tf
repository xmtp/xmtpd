resource "argocd_application" "jaeger" {
  count      = var.enable_monitoring ? 1 : 0
  depends_on = [argocd_project.tools]
  wait       = var.wait_for_ready
  metadata {
    name      = "jaeger"
    namespace = var.argocd_namespace
  }
  spec {
    project = argocd_project.tools.metadata[0].name
    source {
      repo_url        = "https://jaegertracing.github.io/helm-charts"
      chart           = "jaeger"
      target_revision = "0.67.6"
      helm {
        release_name = "jaeger"
        values       = <<EOT
          allInOne:
            enabled: true
            args:
              - --memory.max-traces=10000
              - --query.base-path=/jaeger/ui
              - --prometheus.server-url=http://prometheus-server:80
            extraEnv:
              - name: COLLECTOR_OTLP_ENABLED
                value: "true"
              - name: METRICS_STORAGE_TYPE
                value: prometheus
            nodeSelector:
              node-pool: ${var.node_pool}
            ingress:
              enabled: true
              ingressClassName: ${var.ingress_class_name}
              hosts:
                - jaeger.localhost
          provisionDataStore:
            cassandra: false
          agent:
            enabled: false
          collector:
            enabled: false
          query:
            enabled: false
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
