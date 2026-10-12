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

# Replace the cluster_id with a real KES cluster ID in your account.
data "ksyun_kes_instance_groups" "default" {
  cluster_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  output_file = "out_file"
}

output "first_instance_group_id" {
  value = data.ksyun_kes_instance_groups.default.instance_groups.0.id
}
