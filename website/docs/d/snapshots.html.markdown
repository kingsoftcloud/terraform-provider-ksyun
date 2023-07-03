---
subcategory: "Volume(EBS)"
layout: "ksyun"
page_title: "ksyun: ksyun_snapshots"
sidebar_current: "docs-ksyun-datasource-snapshots"
description: |-
  This data source provides a list of EBS snapshots.
---

# ksyun_snapshots

This data source provides a list of EBS snapshots.

#

## Example Usage

```hcl
data "ksyun_snapshots" "default" {
  output_file = "output_result"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) availability zone.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `snapshot_id` - (Optional) The Id of the snapshot.
* `snapshot_name` - (Optional) The name of the snapshot.
* `volume_category` - (Optional) The category of the volume.
* `volume_id` - (Optional) The ID of the volume.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `snapshots` - An information list of EBS snapshot. Each element contains the following attributes:
  * `availability_zone` - availability zone.
  * `create_time` - creation time.
  * `progress` - Snapshot progress. Example value: 100%.
  * `size` - snapshot size, unit: GB.
  * `snapshot_id` - The Id of the snapshot.
  * `snapshot_name` - The name of the snapshot.
  * `snapshot_status` - snapshot status.
  * `snapshot_type` - snapshot type.
  * `volume_category` - The category of the volume.
  * `volume_id` - The Id of the volume.
  * `volume_status` - Volume status.


