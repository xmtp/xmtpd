variable "name" {}
variable "kubeconfig_path" {}
variable "nodes" {
  type = list(object({
    labels              = optional(map(string))
    extra_port_mappings = optional(map(number))
  }))
}
variable "wait_for_ready" { default = true }
