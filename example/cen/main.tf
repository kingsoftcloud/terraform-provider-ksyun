# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

resource "ksyun_cen" "default" {
  cen_name="cen1"
  description="zice"
}
