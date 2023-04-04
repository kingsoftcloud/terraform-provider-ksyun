---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_security_groups"
sidebar_current: "docs-ksyun-datasource-security_groups"
description: |-
  This data source provides a list of Security Group resources according to their Security Group ID, name and resource id.
---

# ksyun_security_groups

This data source provides a list of Security Group resources according to their Security Group ID, name and resource id.

#

## Example Usage

```hcl
data "ksyun_security_groups" "default" {
  output_file = "output_result"
  ids         = []
  vpc_id      = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Security Group IDs, all the Security Group resources belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `vpc_id` - (Optional) A list of VPC IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_groups` - An information list of security groups. Each element contains the following attributes:
  * `create_time` - The time of creation for the security group.
  * `id` - The ID of the security group.
  * `name` - The name of the security group.
  * `security_group_entry_set` - A list of the security group entries.
    * `cidr_block` - The cidr block of source.
    * `description` - The description of the security group entry.
    * `direction` - The direction of the entry.
    * `icmp_code` - ICMP code.
    * `icmp_type` - ICMP type.
    * `port_range_from` - The start of port numbers.
    * `port_range_to` - The end of port numbers.
    * `protocol` - protocol of the entry.
    * `security_group_entry_id` - The ID of the security group entry.
  * `security_group_id` - The ID of the security group.
  * `security_group_name` - The name of the security group.
  * `security_group_type` - The type of the security group.
  * `vpc_id` - The ID of the VPC.
* `total_count` - Total number of Security Group resources that satisfy the condition.


