variable "kubeconfig_path" { default = ".xmtp/kubeconfig.yaml" }
variable "node_container_image" { default = "xmtpdev/xmtpd:dev" }
variable "nodes" {
  type = list(object({
    name                 = string
    node_id              = string
    p2p_public_address   = string
    p2p_persistent_peers = list(string)
  }))
}
variable "node_keys" {
  type      = map(string)
  sensitive = true
}
variable "enable_chat_app" { default = true }
variable "enable_monitoring" { default = true }
variable "num_xmtp_node_pool_nodes" { default = 2 }
