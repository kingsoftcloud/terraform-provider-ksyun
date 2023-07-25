# Specify the provider and access details
provider "ksyun" {
}

# Get  dnats
data "ksyun_dnats" "default" {
  private_ip_address = "10.7.x.xxx"
  nat_id             = "5c7b7925-xxxx-xxxx-xxxx-434fc8042329"
  dnat_ids           = ["5cxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx"]
  output_file        = "output_result"
}
