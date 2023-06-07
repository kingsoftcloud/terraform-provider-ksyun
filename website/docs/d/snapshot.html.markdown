---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_snapshot"
sidebar_current: "docs-ksyun-datasource-snapshot"
description: |-
  Query instance auto snapshot policies information
---

# ksyun_snapshot

Query instance auto snapshot policies information

#

## Example Usage

```hcl
data "ksyun_snapshot" "foo" {
  name                     = "your auto snapshot policy name"
  auto_snapshot_policy_ids = ["auto snapshot policy id"] // a list of auto snapshot policy id that can be null
  output_file              = "output_result_snapshot"
}

output "ksyun_snapshot" {
  value = data.ksyun_snapshot.foo
}

output "ksyun_snapshots_total_count" {
  value = data.ksyun_snapshot.foo.total_count
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) the name of KEC snapshot policy.
* `auto_snapshot_policy_ids` - (Optional) The id of auto snapshot policy.
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
* `total_count` - Total number of snapshot policies resources that satisfy the condition.


