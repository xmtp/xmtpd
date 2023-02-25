variable "name" { default = "chat-app" }
variable "namespace" {}
variable "node_pool_label_key" {}
variable "node_pool_label_value" {}
variable "service_port" { default = 80 }
variable "container_port" { default = 3000 }
variable "container_image" { default = "snormorexmtp/chat-app:latest" }
variable "api_url" {}
variable "hostnames" { type = list(string) }
variable "ingress_class_name" {}
