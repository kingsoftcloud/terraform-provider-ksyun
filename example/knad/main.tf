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
# Create an knad
resource "ksyun_knad" "default" {
  link_type = "DDoS_BGP"
  ip_count = 10
  band = 30
  max_band = 80
  idc_band = 150
  duration = 1
  knad_name = "test-tf5"
  bill_type = 5
  service_id = "KEAD_30G"
  project_id="0"
}