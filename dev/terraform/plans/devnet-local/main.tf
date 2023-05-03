module "cluster" {
  source = "git@github.com:xmtp-labs/xmtpd-terraform.git//modules/xmtp-cluster-kind?ref=12d1e46"

  # Uncomment this line and comment out the previous source line to use a
  # local instance of xmtpd-modules living in the parent directory of xmtpd.
  # source = "../../../../../xmtpd-terraform/modules/xmtp-cluster-kind"

  name_prefix                 = "xmtp-devnet"
  nodes                       = var.nodes
  node_keys                   = var.node_keys
  node_container_image        = var.node_container_image
  e2e_container_image         = var.e2e_container_image
  e2e_delay                   = var.e2e_delay
  enable_chat_app             = var.enable_chat_app
  enable_e2e                  = var.enable_e2e
  enable_monitoring           = var.enable_monitoring
  node_container_cpu_limit    = var.node_container_cpu_limit
  node_container_memory_limit = var.node_container_memory_limit

  ingress_http_port  = 80
  ingress_https_port = 443
}
