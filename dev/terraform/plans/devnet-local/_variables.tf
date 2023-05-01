variable "node_container_image" { default = "xmtpdev/xmtpd:dev" }
variable "nodes" {
  type = list(object({
    name                 = string
    node_id              = string
    p2p_public_address   = string
    p2p_persistent_peers = list(string)
    store_type           = optional(string, "mem")
  }))
}
variable "node_keys" {
  type      = map(string)
  sensitive = true
}
variable "enable_chat_app" { default = true }
variable "enable_monitoring" { default = true }
variable "e2e_delay" { default = "" }
variable "node_container_cpu_limit" { default = "500m" }
variable "node_container_memory_limit" { default = "500Mi" }
