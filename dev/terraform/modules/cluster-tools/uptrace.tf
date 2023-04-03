locals {
    uptrace_endpoint = "uptrace:14317"
    uptrace_hosts = [
      for name in var.uptrace_hostnames:
      { host: name, paths: [ { path: "/", pathType: "Prefix" } ] }
    ]
}

module "argocd_app_uptrace" {
  count  = var.enable_monitoring ? 1 : 0
  source = "../argocd-application"

  argocd_namespace = var.argocd_namespace
  argocd_project   = module.argocd_project.name
  name             = "uptrace"
  namespace        = var.namespace
  wait             = var.wait_for_ready
  repo_url         = "https://charts.uptrace.dev"
  chart            = "uptrace"
  target_revision  = "1.3.1"
  helm_values = [
    <<EOF
      ingress:
        enabled: true
        hosts: ${jsonencode(local.uptrace_hosts)}
      uptrace:
        config:
          projects:
            - id: 1
              name: uptrace
              token: 12345
              pinned_attrs:
                - service
                - host.name
                - deployment.environment
            - id: 2
              name: devnet
              token: 6789
              pinned_attrs:
                - service
                - host.name
          site:
            addr: 'http:uptrace.localhost'
    EOF
  ]
}
