---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_auto_snapshot_policy"
sidebar_current: "docs-ksyun-datasource-auto_snapshot_policy"
description: |-
  Query instance auto snapshot policies information
---

# ksyun_auto_snapshot_policy

Query instance auto snapshot policies information

#

## Example Usage

## query auto snapshot policy with name or id

```hcl
data "ksyun_auto_snapshot_policy" "foo" {
  name                     = "your auto snapshot policy name"
  auto_snapshot_policy_ids = ["auto snapshot policy id"]
  output_file              = "output_result_snapshot"
}

output "ksyun_auto_snapshot_policy" {
  value = data.ksyun_auto_snapshot_policy.foo
}
```

## query all auto snapshot policy

```hcl
data "ksyun_auto_snapshot_policy" "foo" {
  output_file = "output_result_snapshot"
}
```

## Argument Reference

The following arguments are supported:

* `auto_snapshot_policy_ids` - (Optional) The id of auto snapshot policy.
* `name` - (Optional) the name of auto snapshot policy.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `snapshots` - An information list of auto snapshot policy. Each element contains the following attributes:
  * `attach_ebs_volume_num` - The snapshot retention period (unit: day).
  * `attach_local_volume_num` - The volume number that is attached to this policy.
  * `auto_snapshot_date` - The snapshot policy will be triggered in these dates per month.
  * `auto_snapshot_policy_id` - The snapshot policy id.
  * `auto_snapshot_policy_name` - The snapshot policy name.
  * `auto_snapshot_time` - The snapshot policy will be created in these hours.
  * `creation_date` - The snapshot policy creation date.
* `total_count` - Total number of auto snapshot policies resources that satisfy the condition.


