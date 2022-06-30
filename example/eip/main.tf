# Specify the provider and access details
//provider "ksyun" {
//  region = "cn-beijing-6"
//}
//data "ksyun_lines" "default" {
//  output_file="output_result1"
//  line_name="BGP"
//}

# Create an eip
resource "ksyun_eip" "default1" {
  band_width =1
  charge_type = "Daily"
  project_id = 0
  tags = {
    "xym-test" ="123",
  }
}
