# Specify the provider and access details
terraform {
  required_providers {
    ksyun = {
        source = "kingsoftcloud/ksyun"
    }
  }
}

provider "ksyun" {
  region = "cn-beijing-6"

}

# Get  knads
data "ksyun_knads" "default" {

  output_file="output_result"

  project_id=["0"]
   service_id= "KEAD_30G"
    band=30
    max_band=30
    ip_count= 10
    bill_type= 1
    idc_band= 100
    knad_name = "terraform1"
}

