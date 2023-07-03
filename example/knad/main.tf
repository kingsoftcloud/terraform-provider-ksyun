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
# Create an eip
resource "ksyun_knad" "default2" {
  link_type = "DDoS_BGP"
  ip_count = 10
  band = 30
  max_band = 80
  idc_band = 150
  duration = 1
  knad_name = "test-tf2"
  bill_type = 1
  service_id = "KEAD_30G"
  project_id="123"
}