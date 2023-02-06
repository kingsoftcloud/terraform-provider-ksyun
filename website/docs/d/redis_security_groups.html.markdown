---
subcategory: "Redis"
layout: "ksyun"
page_title: "ksyun: ksyun_redis_security_groups"
sidebar_current: "docs-ksyun-datasource-redis_security_groups"
description: |-
  Provides a list of Redis security groups in the current region.
---

# ksyun_redis_security_groups

Provides a list of Redis security groups in the current region.

#

## Example Usage

```hcl
data "ksyun_redis_security_groups" "default" {
  output_file = "output_result1"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - An information list of Redis security groups. Each element contains the following attributes:
  * `created` - creation time.
  * `description` - security group description.
  * `name` - security group name.
  * `security_group_id` - security group ID.
  * `updated` - updated time.
* `total_count` - Total number of Redis security groups that satisfy the condition.


