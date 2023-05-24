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

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `local_snapshot_set` - An information list of KEC local snapshots. Each element contains the following attributes:
* `total_count` - Total number of local snapshots that satisfy the condition.


