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

# Create a KES8 (vector) cluster. KES8 clusters run Elasticsearch 8.x / 9.x and
# require Xpack credentials. Replace the VPC / security group ids with real
# resources in your account before applying. The cluster is created via the
# shared kmr/2021-09-02 LaunchCluster backend switched into KES mode by
# ClusterType="KES" (handled internally by the provider); on that backend the
# per-group placement (availability_zone / vpc_id / vpc_subnet_id) and the
# per-group index / charge_type live INSIDE each instance_groups entry.
resource "ksyun_kes_cluster_vector" "example" {
  cluster_name      = "my-kes-vector-cluster"
  main_version      = "8.19.13"
  distribution      = "ElasticSearch"
  scenario          = "General"
  charge_type       = "Monthly"
  project_id        = 0
  vpc_domain_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  vpc_subnet_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  security_group_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

  # required for any 8.x / 9.x cluster
  xpack_username = "elastic"
  xpack_password = "your-xpack-password"

  # KES8-specific switches (defaults shown).
  xpack_enabled             = true
  enable_best_delivery_mode = false
  enable_eip                = false
  instance_login_type       = "later"
  services                  = ["elasticsearch"]
  plugins                   = ["repository-ks3"]

  instance_groups {
    instance_group_type     = "MASTER"
    instance_type           = "ES.I3.2C4G"
    instance_type_code      = "ES.I3.2C4G"
    resource_type           = "KEC"
    instance_count          = 3
    instance_group_index    = 1
    availability_zone       = "cn-beijing-6a"
    availability_zone_index = 1
    vpc_id                  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    vpc_subnet_id           = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    charge_type             = "Monthly"
    volume_type             = "SSD3.0"
    volume_size             = 100
    volume_count            = 1
    system_disk_type        = "SSD3.0"
    system_disk_size        = 50
  }

  instance_groups {
    instance_group_type     = "DATA"
    instance_type           = "ES.D4.8C16G"
    instance_type_code      = "ES.D4.8C16G"
    resource_type           = "KEC"
    instance_count          = 3
    instance_group_index    = 2
    availability_zone       = "cn-beijing-6a"
    availability_zone_index = 1
    vpc_id                  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    vpc_subnet_id           = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    charge_type             = "Monthly"
    volume_type             = "SSD3.0"
    volume_size             = 500
    volume_count            = 1
    system_disk_type        = "SSD3.0"
    system_disk_size        = 50
  }
}

output "cluster_id" {
  value = ksyun_kes_cluster_vector.example.cluster_id
}
