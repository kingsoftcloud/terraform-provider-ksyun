# Specify the provider and access details
terraform {
  required_providers {
    ksyun = {
      source = "kingsoftcloud/ksyun"
    }
  }
}

# Create an knadAssociate
resource "ksyun_knad_associate" "default" {
  knad_id="xxxxxx-xxxxf-xxxxxx"
  ip = ["1.1.1.1"]
}