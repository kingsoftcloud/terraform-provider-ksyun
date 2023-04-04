---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_subnet_allocated_ip_addresses"
sidebar_current: "docs-ksyun-datasource-subnet_allocated_ip_addresses"
description: |-
  This data source provides a list of allocated IPs.
---

# ksyun_subnet_allocated_ip_addresses

This data source provides a list of allocated IPs.

#

## Example Usage

```hcl
data "ksyun_subnet_allocated_ip_addresses" "default" {
  output_file = "output_result"
  ids         = ["494c3a64-eff9-4438-aa7c-694b7baxxxxx"]
  subnet_id   = ["494c3a64-eff9-4438-aa7c-694b7baxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of subnet IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `subnet_id` - (Optional) The ID of the subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `subnet_allocated_ip_addresses` - A list of allocated IPs.
* `total_count` - Total number of Allocated IPs that satisfy the condition.


