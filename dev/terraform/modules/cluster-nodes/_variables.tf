variable "namespace" {}
variable "container_image" {}
variable "node_pool_label_key" {}
variable "node_pool" {}
variable "ingress_class_name" {}
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
variable "wait_for_ready" {}
variable "hostnames" { type = list(string) }
variable "node_api_http_port" { type = number }
variable "storage_class_name" {}
variable "container_storage_request" {}
variable "container_cpu_request" {}
variable "one_instance_per_k8s_node" { type = bool }
variable "debug" { type = bool }
variable "argocd_project" {}
variable "argocd_namespace" {}
