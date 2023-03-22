variable "name" {}
variable "namespace" {}
variable "description" { default = "" }
variable "source_repos" { default = ["*"] }
variable "destinations" {
  type = list(object({ server : string, namespace : string }))
}
variable "cluster_resource_whitelist" {
  type = list(object({ kind : string, group : string }))
  default = [{
    kind  = "*"
    group = "*"
  }]
}
variable "namespace_resource_whitelist" {
  type = list(object({ kind : string, group : string }))
  default = [{
    kind  = "*"
    group = "*"
  }]
}
variable "namespace_resource_blacklist" {
  type    = list(object({ kind : string, group : string }))
  default = null
}
