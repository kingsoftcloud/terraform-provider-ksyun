---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_snapshot_volume_association"
sidebar_current: "docs-ksyun-resource-snapshot_volume_association"
description: |-
  Provides a tag resource.
---

# ksyun_snapshot_volume_association

Provides a tag resource.

#

## Example Usage

```hcl
resource "ksyun_snapshot_volume_association" "foo" {
  attach_volume_id        = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  auto_snapshot_policy_id = "auto_snapshot_policy_id"
}
```

## Argument Reference

The following arguments are supported:

* `attach_volume_id` - (Required, ForceNew) The id of the volume.
* `auto_snapshot_policy_id` - (Required, ForceNew) The id of the auto_snapshot_policy_id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



