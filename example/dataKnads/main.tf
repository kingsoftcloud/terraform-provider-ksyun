# Specify the provider and access details
terraform {
  required_providers {
    ksyun = {
      source = "kingsoftcloud/ksyun"
    }
  }
}


# Get  knads
data "ksyun_knads" "default" {
  output_file = "output_result"
  project_id = []
  ids    = []
}