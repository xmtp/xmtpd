variable "name" { default = "promlens" }
variable "namespace" {}
variable "node_pool_label_key" {}
variable "node_pool_label_value" {}
variable "service_port" { default = 80 }
variable "container_port" { default = 4000 }
variable "container_image" { default = "prom/promlens:latest" }
variable "hostnames" { type = list(string) }
variable "ingress_class_name" {}
variable "prometheus_server_endpoint" {}
