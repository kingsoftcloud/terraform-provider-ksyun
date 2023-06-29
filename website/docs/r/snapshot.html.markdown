---
subcategory: "Volume(EBS)"
layout: "ksyun"
page_title: "ksyun: ksyun_snapshot"
sidebar_current: "docs-ksyun-resource-snapshot"
description: |-
  Provides a EBS snapshot resource.
---

# ksyun_snapshot

Provides a EBS snapshot resource.

#

## Example Usage

```hcl
resource "ksyun_snapshot" "default" {
  snapshot_name = "test_tf_snapshot"
  snapshot_desc = "test descrition"
  volume_id     = "xxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `volume_id` - (Required, ForceNew) The ID of the volume. Snapshot requires the Volume to be in "in-use" or "available" status.When the Volume status is "in-use", the kec instance status can be either "running" or "stopped".
* `snapshot_desc` - (Optional) The description of the snapshot.
* `snapshot_name` - (Optional) The name of the snapshot.
* `snapshot_type` - (Optional, ForceNew) The type of the snapshot, valid values: 'LocalSnapShot', 'CommonSnapShot'. Default is 'CommonSnapShot'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `availability_zone` - Availability zone.
* `create_time` - The creation time.
* `progress` - Snapshot progress. Example value: 100%.
* `size` - The size of the snapshot, unit is 'GB'.
* `snapshot_status` - snapshot status.
* `volume_category` - The category of the volume, 'data' or 'system'.
* `volume_status` - Volume status.


## Import

Instance can be imported using the `id`, e.g.

```
$ terraform import ksyun_snapshot.default xxxxxx
```

