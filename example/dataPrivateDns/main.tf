data "ksyun_private_dns_zones" "foo" {
  output_file = "pdns_output_result"
  zone_ids = []
}

data "ksyun_private_dns_records" "foo" {
  zone_id = "" // Required
  output_file = "pdns_records_output_result"
  region_name = ["cn-beijing-6"]
  record_ids = []
}