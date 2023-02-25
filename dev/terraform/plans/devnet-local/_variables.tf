variable "kubeconfig_path" { default = ".xmtp/kubeconfig.yaml" }
variable "node_container_image" { default = "xmtp/xmtpd:dev" }
variable "e2e_container_image" { default = "xmtp/xmtpd-e2e:dev" }
variable "nodes" {
  type = list(object({
    name               = string
    node_id            = string
    p2p_public_address = string
    persistent_peers   = list(string)
  }))
}
variable "node_keys" { type = map(string) }
variable "enable_e2e" { default = true }
variable "enable_demo_app" { default = true }
variable "enable_monitoring" { default = true }
variable "num_xmtp_node_pool_nodes" { default = 2 }
