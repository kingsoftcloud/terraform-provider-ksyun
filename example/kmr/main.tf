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
  endpoint = "kmr.api.ksyun.com"
}

data "ksyun_kmr_clusters" "default" {
  marker      = "limit=10&offset=0"
  output_file = "out_file"
}

output "kmr_clusters" {
  value = data.ksyun_kmr_clusters.default.clusters
}

output "total" {
  value = data.ksyun_kmr_clusters.default.total
}
