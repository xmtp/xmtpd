locals {
  jaeger_hostnames          = [for hostname in var.hostnames : "jaeger.${hostname}"]
  jaeger_public_hostname    = local.jaeger_hostnames[0]
  jaegar_public_url         = "http://${local.jaeger_public_hostname}"
  jaeger_collector_endpoint = "jaeger-collector:4317"
  jaeger_query_endpoint     = "jaeger-query:16686"
}

module "argocd_app_jaeger" {
  count  = var.enable_monitoring ? 1 : 0
  source = "../argocd-application"

  argocd_namespace = var.argocd_namespace
  argocd_project   = module.argocd_project.name
  name             = "jaeger"
  namespace        = var.namespace
  wait             = var.wait_for_ready
  repo_url         = "https://jaegertracing.github.io/helm-charts"
  chart            = "jaeger"
  target_revision  = "0.67.6"
  helm_values = [
    <<EOF
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
    EOF
  ]
}
