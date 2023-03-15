locals {
  jaeger_hostnames          = [for hostname in var.hostnames : "jaeger.${hostname}"]
  jaeger_public_hostname    = local.jaeger_hostnames[0]
  jaegar_public_url         = "http://${local.jaeger_public_hostname}"
  jaeger_collector_endpoint = "jaeger-collector:4317"
  jaeger_query_endpoint     = "jaeger-query:16686"
}

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
              - --prometheus.server-url=${local.prometheus_server_url}
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
                - ${local.jaeger_public_hostname}
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
