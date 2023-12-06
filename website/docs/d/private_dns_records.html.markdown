---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_private_dns_records"
sidebar_current: "docs-ksyun-datasource-private_dns_records"
description: |-
  This data source provides a list of Private Dns Record resources according to their Zone ID.
---

# ksyun_private_dns_records

This data source provides a list of Private Dns Record resources according to their Zone ID.

#

## Example Usage

```hcl
data "ksyun_private_dns_records" "default" {
  output_file = "pdns_records_output_result"
  zone_id     = "a5ae6bf0-0ff4-xxxxxx-xxxxx-xxxxxxxxxx"
  region_name = ["cn-beijing-6"]
  record_ids  = []
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required) Id of the private dns zone. Required.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `record_ids` - (Optional) A list of Record IDs, the Records belong to this private-dns-zone. The value of id is not be `""`.
* `region_name` - (Optional) A list of the filter values that is region name. Such `cn-beijing-6`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `records` - An information list of `Private Dns Records`. Each element contains the following attributes:
  * `create_time` - Creation time.
  * `record_data_set` - The record value and other information like priority and weight etc.
    * `port` - The port of record.
    * `priority` - The priority of record.
    * `record_value` - Record value.
    * `weight` - The weight of record.
  * `record_id` - ID of the record.
  * `record_name` - The name of record.
  * `record_ttl` - The record ttl.
* `total_count` - Total number of Private Dns Record that satisfy the condition.


