# Specify the provider and access details
terraform {
  required_providers {
    ksyun = {
      source = "kingsoftcloud/ksyun"
    }
  }
}

provider "ksyun" {
  #region = "cn-beijing-6"
  region = "cn-qingyangtest-1"

}

# Get  knads
data "ksyun_knads" "default" {
  output_file = "output_result"
  project_id = ["0"]
  ids    = ["knad8b4797ea-3110-3812-9357-f4b7f31fe1a1","knad1c9cabb2-9cd5-36e9-8fee-d4957b9f58bf"]
}

