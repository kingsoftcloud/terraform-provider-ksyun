# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

# Get  knads
data "ksyun_knads" "default" {
  output_file="output_result"

  ids=[]
  project_id=["0"]
  instance_type=["Slb"]
  ip_version = "all"
  network_interface_id=[]
  internet_gateway_id=[]
  band_width_share_id=[]
  line_id=[]
  public_ip=[]
}

