---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_auto_snapshot_policy"
sidebar_current: "docs-ksyun-resource-auto_snapshot_policy"
description: |-
  Provides an auto snapshot policy resource.
---

# ksyun_auto_snapshot_policy

Provides an auto snapshot policy resource.

#

## Example Usage

```hcl
resource "ksyun_auto_snapshot_policy" "foo" {
  name               = "your auto snapshot policy name"
  auto_snapshot_date = [1, 3, 4, 5]
  auto_snapshot_time = [1, 3, 4, 5, 9, 22]
}
```

## Argument Reference

The following arguments are supported:

* `auto_snapshot_date` - (Required) Setting the snapshot date in a week, its scope is between 1 and 7.
* `auto_snapshot_time` - (Required) Setting the snapshot time in a day, its scope is between 0 and 23.
* `name` - (Required) the name of auto snapshot policy.
* `retention_time` - (Optional) the snapshot will be reserved for when, the cap is 9999.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `auto_snapshot_policy_id` - The id of auto snapshot policy.
* `creation_date` - The snapshot policy creation date.


## Import

`ksyun_auto_snapshot_policy` can be imported using the `id`, e.g.

```
$ terraform import ksyun_auto_snapshot_policy.foo "auto_snapshot_policy_id"
```

