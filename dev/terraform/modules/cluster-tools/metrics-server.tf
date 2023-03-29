module "argocd_app_metrics_server" {
  count  = var.enable_monitoring ? 1 : 0
  source = "../argocd-application"

  argocd_namespace = var.argocd_namespace
  argocd_project   = module.argocd_project.name
  name             = "metrics-server"
  namespace        = "kube-system"
  wait             = var.wait_for_ready
  repo_url         = "https://kubernetes-sigs.github.io/metrics-server/"
  chart            = "metrics-server"
  target_revision  = "3.8.3"
  helm_values = [
    <<EOF
      args:
        - --kubelet-insecure-tls
        - --kubelet-preferred-address-types=InternalIP
    EOF
  ]
}
