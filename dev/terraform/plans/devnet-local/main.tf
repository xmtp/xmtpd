module "cluster" {
  source = "git@github.com:xmtp-labs/xmtpd-terraform.git//modules/xmtp-cluster-kind?ref=611159d"

  name_prefix          = "xmtp-devnet"
  nodes                = var.nodes
  node_keys            = var.node_keys
  node_container_image = var.node_container_image
  enable_chat_app      = var.enable_chat_app
  enable_monitoring    = var.enable_monitoring

  ingress_http_port  = 80
  ingress_https_port = 443
}
