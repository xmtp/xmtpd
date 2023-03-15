output "region" {
  value       = var.region
  description = "AWS region"
}

output "public_subnet_cidrs" {
  value       = module.cluster.public_subnet_cidrs
  description = "Public subnet CIDRs"
}

output "private_subnet_cidrs" {
  value       = module.cluster.private_subnet_cidrs
  description = "Private subnet CIDRs"
}

output "vpc_cidr" {
  value       = module.cluster.vpc_cidr
  description = "VPC ID"
}

output "eks_cluster_id" {
  description = "The name of the cluster"
  value       = module.cluster.eks_cluster_id
}

output "eks_cluster_arn" {
  description = "The Amazon Resource Name (ARN) of the cluster"
  value       = module.cluster.eks_cluster_arn
}

output "eks_cluster_endpoint" {
  description = "The endpoint for the Kubernetes API server"
  value       = module.cluster.eks_cluster_endpoint
}

output "eks_cluster_version" {
  description = "The Kubernetes server version of the cluster"
  value       = module.cluster.eks_cluster_version
}

output "eks_cluster_identity_oidc_issuer" {
  description = "The OIDC Identity issuer for the cluster"
  value       = module.cluster.eks_cluster_identity_oidc_issuer
}

output "eks_cluster_managed_security_group_id" {
  description = "Security Group ID that was created by EKS for the cluster. EKS creates a Security Group and applies it to ENI that is attached to EKS Control Plane master nodes and to any managed workloads"
  value       = module.cluster.eks_cluster_managed_security_group_id
}

output "ecr_node_repo_id" {
  value = module.ecr_node_repo.registry_id
}

output "ecr_node_repo_url" {
  value = module.ecr_node_repo.repository_url
}
