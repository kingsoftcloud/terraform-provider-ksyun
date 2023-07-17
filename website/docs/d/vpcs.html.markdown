---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_vpcs"
sidebar_current: "docs-ksyun-datasource-vpcs"
description: |-
  This data source provides a list of VPC resources according to their VPC ID, name.
---

# ksyun_vpcs

This data source provides a list of VPC resources according to their VPC ID, name.

#

## Example Usage

```hcl
data "ksyun_vpcs" "default" {
  output_file = "output_result"
  ids         = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of VPC IDs, all the VPC resources belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter results by VPC name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total_count` - Total number of VPC resources that satisfy the condition.
* `vpcs` - It is a nested type which documented below.
  * `cidr_block` - The CIDR blocks of VPC.
  * `create_time` - The time of creation for VPC.
  * `id` - The ID of VPC.
  * `ipv6_cidr_block_association_set` - An Ipv6 association list of this vpc.
    * `ipv6_cidr_block` - the Ipv6 of this vpc bound.
  * `name` - The name of VPC.
  * `vpc_id` - The ID of VPC.
  * `vpc_name` - The name of VPC.


