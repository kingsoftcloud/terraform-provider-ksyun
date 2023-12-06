---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_private_dns_zones"
sidebar_current: "docs-ksyun-datasource-private_dns_zones"
description: |-
  This data source provides a list of Private Dns Zone.
---

# ksyun_private_dns_zones

This data source provides a list of Private Dns Zone.

#

## Example Usage

```hcl
data "ksyun_private_dns_zones" "foo" {
  output_file = "pdns_output_result"
  zone_ids    = []
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `zone_ids` - (Optional) A list of the filter values that is zone id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total_count` - Total number of Private Dns Zone that satisfy the condition.
* `zones` - An information list of Private Dns Zone. Each element contains the following attributes:
  * `bind_vpc_set` - The zone Bound VPCs.
    * `region_name` - Region Name.
    * `status` - The status of Zone.
    * `vpc_id` - VPC id.
    * `vpc_name` - The VPC name.
  * `create_time` - Creation time.
  * `project_id` - Project Id.
  * `zone_id` - ID of the Private Dns Zone.
  * `zone_name` - The name of Private Dns Zone.
  * `zone_ttl` - Zone TTL.


