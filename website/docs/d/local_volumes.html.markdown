---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_local_volumes"
sidebar_current: "docs-ksyun-datasource-local_volumes"
description: |-
  This data source provides a list of kec local volumes in the current region.
---

# ksyun_local_volumes

This data source provides a list of kec local volumes in the current region.

#

## Example Usage

```hcl
data "ksyun_local_volumes" "default" {
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
  * `creation_date` - creation date.
  * `instance_id` - The ID of the instance.
  * `instance_name` - The name of the instance.
  * `instance_state` - The state of the instance.
  * `local_volume_category` - The category of the local volume.
  * `local_volume_id` - The ID of the local volume.
  * `local_volume_name` - The name of the local volume.
  * `local_volume_size` - The size of the local volume.
  * `local_volume_state` - The state of the local volume.
  * `local_volume_type` - The type of the local volume.
* `total_count` - Total number of local volumes that satisfy the condition.


