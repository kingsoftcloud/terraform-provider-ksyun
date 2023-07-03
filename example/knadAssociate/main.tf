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
  knad_id="knade0e77118-86a9-3761-9c39-b30dc2c0477e"
  ip = ["88.88.88.46"]
}