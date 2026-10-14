terraform {
  required_providers {
    ksyun = {
      source = "kingsoftcloud/ksyun"
    }
  }
}

provider "ksyun" {
  access_key = "your ak"
  secret_key = "your sk"
  region     = "cn-beijing-6"
  endpoint   = "kes.api.ksyun.com"
}

# Create a KES7 (standard) cluster. Replace the VPC / security group ids with
# real resources in your account before applying.
resource "ksyun_kes_cluster" "example" {
  cluster_name      = "my-kes-cluster"
  main_version      = "7.10.0"
  scenario          = "general"
  availability_zone = "cn-beijing-6a"
  charge_type       = "Monthly"
  purchase_time     = 1
  project_id        = 0
  vpc_domain_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  vpc_subnet_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  security_group_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

  instance_groups {
    instance_group_type = "MASTER"
    instance_type       = "C3.2A"
    resource_type       = "host"
    instance_count      = 3
    volume_type         = "SSD3.0"
    volume_size         = 100
  }

  instance_groups {
    instance_group_type = "DATA"
    instance_type       = "C3.4A"
    resource_type       = "host"
    instance_count      = 3
    volume_type         = "SSD3.0"
    volume_size         = 500
  }
}

output "cluster_id" {
  value = ksyun_kes_cluster.example.cluster_id
}
