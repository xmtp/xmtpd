locals {
  ingress_public_hostname = module.system.ingress_public_hostname
}

resource "cloudflare_record" "hostnames" {
  count   = length(var.hostnames)
  zone_id = var.cloudflare_zone_id
  name    = var.hostnames[count.index]
  value   = local.ingress_public_hostname
  type    = "CNAME"
  proxied = true
}

resource "cloudflare_record" "node_hostnames" {
  count   = length(local.node_hostnames)
  zone_id = var.cloudflare_zone_id
  name    = local.node_hostnames[count.index]
  value   = local.ingress_public_hostname
  type    = "CNAME"
  proxied = true
}

resource "cloudflare_record" "argocd_hostnames" {
  count   = length(local.argocd_hostnames)
  zone_id = var.cloudflare_zone_id
  name    = local.argocd_hostnames[count.index]
  value   = local.ingress_public_hostname
  type    = "CNAME"
  proxied = true
}

resource "cloudflare_record" "grafana_hostnames" {
  count   = length(local.grafana_hostnames)
  zone_id = var.cloudflare_zone_id
  name    = local.grafana_hostnames[count.index]
  value   = local.ingress_public_hostname
  type    = "CNAME"
  proxied = true
}

resource "cloudflare_record" "jaeger_hostnames" {
  count   = length(local.jaeger_hostnames)
  zone_id = var.cloudflare_zone_id
  name    = local.jaeger_hostnames[count.index]
  value   = local.ingress_public_hostname
  type    = "CNAME"
  proxied = true
}

resource "cloudflare_record" "prometheus_hostnames" {
  count   = length(local.prometheus_hostnames)
  zone_id = var.cloudflare_zone_id
  name    = local.prometheus_hostnames[count.index]
  value   = local.ingress_public_hostname
  type    = "CNAME"
  proxied = true
}

resource "cloudflare_record" "chat_app_hostnames" {
  count   = length(local.chat_app_hostnames)
  zone_id = var.cloudflare_zone_id
  name    = local.chat_app_hostnames[count.index]
  value   = local.ingress_public_hostname
  type    = "CNAME"
  proxied = true
}

