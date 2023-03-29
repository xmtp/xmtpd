variable "name" {}
variable "argocd_namespace" {}
variable "argocd_project" {}
variable "repo_url" {}
variable "target_revision" {}
variable "chart" {}
variable "path" { default = "" }
variable "helm_values" {
  type    = list(string)
  default = []
}
variable "destination_server" { default = "https://kubernetes.default.svc" }
variable "namespace" {}
variable "auto_sync" { default = true }
variable "wait" { default = false }
