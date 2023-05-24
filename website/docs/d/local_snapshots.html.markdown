---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_local_snapshots"
sidebar_current: "docs-ksyun-datasource-local_snapshots"
description: |-
  This data source provides a list of kec local snapshots in the current region.
---

# ksyun_local_snapshots

This data source provides a list of kec local snapshots in the current region.

#

## Example Usage

```hcl
data "ksyun_local_snapshots" "default" {
  output_file = ""
}
```

## Argument Reference

The following arguments are supported:

* `local_volume_name` - (Optional) The name of the volume.
* `local_volume_snapshot_id` - (Optional) The ID of the snapshot.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `source_local_volume_id` - (Optional) The ID of the volume.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `local_snapshot_set` - An information list of KEC local snapshots. Each element contains the following attributes:
  * `creation_date` - creation date.
  * `disk_size` - Disk size.
  * `instance_id` - The ID of the instance.
  * `local_volume_snapshot_desc` - The description of snapshot.
  * `local_volume_snapshot_id` - The ID of snapshot.
  * `local_volume_snapshot_name` - The name of snapshot.
  * `snapshot_type` - snapshot type.
  * `source_local_volume_category` - The category of the volume.
  * `source_local_volume_id` - The ID of the volume.
  * `source_local_volume_name` - The name of the volume.
  * `source_local_volume_state` - The state of the volume.
  * `state` - The state of the snapshot.
* `total_count` - Total number of local snapshots that satisfy the condition.


