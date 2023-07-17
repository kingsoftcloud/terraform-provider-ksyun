---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_auto_snapshot_volume_association"
sidebar_current: "docs-ksyun-datasource-auto_snapshot_volume_association"
description: |-
  Provides a list of Redis security groups in the current region.
---

# ksyun_auto_snapshot_volume_association

Provides a list of Redis security groups in the current region.

#

## Example Usage

query asp and volume associations by volume id

```hcl
data "ksyun_auto_snapshot_volume_association" "foo" {
  output_file      = "output_result_volume_id"
  attach_volume_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

query all associations.

```hcl
data "ksyun_auto_snapshot_volume_association" "foo1" {
  output_file = "output_result_null"
}
```

query asp and volume associations by auto_snapshot_policy_id

```hcl
data "ksyun_auto_snapshot_volume_association" "foo2" {
  output_file             = "output_result_policy_id"
  auto_snapshot_policy_id = "auto_snapshot_policy_id"
}
```

## Argument Reference

The following arguments are supported:

* `attach_volume_id` - (Optional) The id of the volume.
* `auto_snapshot_policy_id` - (Optional) The id of the auto snapshot policy.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total_count` - Total number of snapshot policies resources that satisfy the condition.
* `volume_asp_associations` - An information list of the associations of volumes and auto snapshot policy. Each element contains the following attributes:
  * `attach_volume_id` - The id of the volume.
  * `auto_snapshot_policy_id` - The snapshot policy id.


