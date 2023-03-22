variable "name" {}
variable "namespace" {}
variable "private_key" { sensitive = true }
variable "p2p_persistent_peers" { type = list(string) }
variable "container_image" {}
variable "hostnames" { type = list(string) }
variable "storage_class_name" {}
variable "storage_request" {}
variable "cpu_request" {}
variable "p2p_port" { type = number }
variable "api_http_port" { type = number }
variable "api_grpc_port" { type = number }
variable "metrics_port" { type = number }
variable "node_pool_label_key" {}
variable "node_pool" {}
variable "one_instance_per_k8s_node" { type = bool }
variable "ingress_class_name" {}
variable "wait_for_ready" { type = bool }
variable "debug" { type = bool }
variable "store_type" {
  type        = string
  description = "type of persistent store to use"
  default     = "mem"
  validation {
    condition     = contains(["mem", "bolt", "postgres"], var.store_type)
    error_message = "Recognized store types are mem, bolt or postgres"
  }
}
variable "argocd_project" {}
variable "argocd_namespace" {}
