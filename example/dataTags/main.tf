# Specify the provider and access details
provider "ksyun" {
}

# Get  tags
data "ksyun_tags" "default" {
  output_file="output_result"

  # optional
  # eg. key = ["tag_key1", "tag_key2", ...]
  keys = []
  # optional
  # eg. value = ["tag_value1", ...]
  values = []
  # optional
  # eg. resource_type = ["kec-instance", "eip", ...]
  resource_types = []
  # optional
  # eg. key = ["instance_uuid", ...]
  resource_ids = []

}

