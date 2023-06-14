---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_data_guard_group"
sidebar_current: "docs-ksyun-datasource-data_guard_group"
description: |-
  Query instance auto snapshot policies information
---

# ksyun_data_guard_group

Query instance auto snapshot policies information

#

## Example Usage

```hcl
data "ksyun_data_guard_group" "foo" {
  output_file = "output_result"
}
data "ksyun_data_guard_group" "foo1" {
  data_guard_name = "Data Guard Name"
}

data "ksyun_data_guard_group" "foo2" {
  data_guard_id = "Data Guard Id"
}
```

## Argument Reference

The following arguments are supported:

* `data_guard_id` - (Optional) The id of data guard group.
* `data_guard_name` - (Optional) The name of data guard group.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data_guard_groups` - An information list of data guard groups. Each element contains the following attributes:
  * `data_guard_capacity` - The capacity of data guard group.
  * `data_guard_id` - The data guard group id.
  * `data_guard_instances_list` - The data guard group includes instances.
  * `data_guard_level` - The data guard group level, if the value is Host represent machine level, and the tol represent the domain of disaster tolerance.
  * `data_guard_name` - The data guard group name.
  * `data_guard_type` - The data guard group display type.
  * `data_guard_used_size` - This data guard group includes the amount of instances.
* `total_count` - Total number of snapshot policies resources that satisfy the condition.


