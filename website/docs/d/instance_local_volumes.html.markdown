---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_instance_local_volumes"
sidebar_current: "docs-ksyun-datasource-instance_local_volumes"
description: |-
  This data source provides a list of kec local volumes in the current region.
---

# ksyun_instance_local_volumes

This data source provides a list of kec local volumes in the current region.

#

## Example Usage

```hcl
data "ksyun_instance_local_volumes" "default" {
  output_file = ""
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Optional) The name of the instance which the volume belong to.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `local_volume_set` - An information list of KEC local volumes. Each element contains the following attributes:
* `total_count` - Total number of local volumes that satisfy the condition.


