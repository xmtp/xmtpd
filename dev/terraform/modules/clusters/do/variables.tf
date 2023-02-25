variable "cluster_name" {}
variable "do_token" { sensitive = true }
variable "cluster_region" {}
variable "node_pool_label_key" {}
variable "num_xmtp_nodes_pool_nodes" { type = number }
variable "num_xmtp_utils_pool_nodes" { type = number }
variable "default_pool_label_value" { default = "default" }
variable "xmtp_nodes_pool_label_value" {}
variable "xmtp_utils_pool_label_value" {}
variable "default_pool_node_size" { default = 1 }
variable "xmtp_nodes_pool_node_size" {}
variable "xmtp_utils_pool_node_size" {}
