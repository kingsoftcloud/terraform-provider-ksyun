---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_nats"
sidebar_current: "docs-ksyun-datasource-nats"
description: |-
  This data source provides a list of Nat resources according to their Nat ID and the VPC they belong to.
---

# ksyun_nats

This data source provides a list of Nat resources according to their Nat ID and the VPC they belong to.

#

## Example Usage

```hcl
data "ksyun_nats" "default" {
  output_file = "output_result"
  ids         = []
  vpc_ids     = []
  project_ids = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Nat IDs, all the Nat resources belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter results by NAT name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `project_ids` - (Optional) A list of Project id that the desired Nat belongs to.
* `vpc_ids` - (Optional) A list of VPC id that the desired Nat belongs to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `nats` - An information list of NATs. Each element contains the following attributes:
  * `associate_nat_set` - The subnet associate list of the desired Nat.
    * `subnet_id` - The ID of the subnet.
  * `band_width` - The nat ip band width of the desired Nat.
  * `create_time` - The time of creation of Nat.
  * `id` - The ID of NAT.
  * `nat_ip_number` - The nat ip count of the desired Nat.
  * `nat_ip_set` - The nat ip list of the desired Nat.
    * `nat_ip_id` - The ID of the NAT IP.
    * `nat_ip` - NAT IP address.
  * `nat_mode` - The mode of the NAT.
  * `nat_name` - The name of the NAT.
  * `nat_type` - The type of the NAT.
  * `project_id` - The ID of the project.
  * `vpc_id` - The VPC ID of the desired Nat belongs to.
* `total_count` - Total number of NAT resources that satisfy the condition.


