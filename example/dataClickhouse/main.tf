provider "ksyun"{
  region = "cn-beijing-6"
}

data "ksyun_clickhouse" "default" {

  output_file = "output_file"

  instance_id = "f20bcde7-c428-43f1-9170-364ebd959e9c"

  product_type = "ClickHouse_Single | ClickHouse"
  project_ids = "107147"
  tag_id = "32960"

  fuzzy_search = "20251105192650"

  offset = 0
  limit = 2
}

