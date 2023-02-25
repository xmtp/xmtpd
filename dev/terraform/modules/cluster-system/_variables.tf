variable "namespace" {}
variable "node_pool_label_key" {}
variable "node_pool" {}
variable "cluster_http_node_port" { type = number }
variable "cluster_https_node_port" { type = number }
variable "argocd_hostnames" { type = list(string) }
variable "ingress_class_name" {}
