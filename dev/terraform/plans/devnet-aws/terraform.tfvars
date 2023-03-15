region                       = "us-east-2"
availability_zones           = ["us-east-2a", "us-east-2b"]
vpc_cidr_block               = "172.16.0.0/16"
kubernetes_version           = "1.25"
enabled_cluster_log_types    = ["audit"]
cluster_log_retention_period = 7

hostnames = ["xmtp.snormore.dev", "xmtp-devnet-aws-us-east-2.snormore.dev"]
