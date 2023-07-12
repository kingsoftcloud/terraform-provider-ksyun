# Specify the provider and access details
terraform {
  required_providers {
    ksyun = {
      source = "kingsoftcloud/ksyun"
    }
  }
}

# Create an knad
resource "ksyun_knad" "default" {
  link_type = "DDoS_BGP"
  ip_count = 10
  band = 30
  max_band = 30
  idc_band = 100
  duration = 1
  knad_name = "terraform1"
  bill_type = 5
  service_id = "KNAD_30G"
  project_id="0"
}