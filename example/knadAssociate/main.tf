# Specify the provider and access details
terraform {
  required_providers {
    ksyun = {
      source = "kingsoftcloud/ksyun"
    }
  }
}

provider "ksyun" {
  region = "cn-qingyangtest-1"
}
# Create an knadAssociate
resource "ksyun_knad_associate" "default" {
  knad_id="knadba4d704f-35b1-3354-8d0f-xxxxxx"
  ip = ["88.88.88.46"]
}