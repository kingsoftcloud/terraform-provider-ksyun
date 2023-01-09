---
subcategory: "Volume(EBS)"
layout: "ksyun"
page_title: "ksyun: ksyun_volumes"
sidebar_current: "docs-ksyun-datasource-volumes"
description: |-
  This data source provides a list of EBS volumes.
---

# ksyun_volumes

This data source provides a list of EBS volumes.

#

## Example Usage

```hcl
data "ksyun_volumes" "default" {
  output_file       = "output_result"
  ids               = []
  volume_category   = ""
  volume_status     = ""
  volume_type       = ""
  availability_zone = ""
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) The availability zone in which the EBS volume resides.
* `ids` - (Optional) A list of EBS IDs, all the EBS resources belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `volume_category` - (Optional) The category to which the EBS volume belongs.
* `volume_create_date` - (Optional) The time when the EBS volume was created.
* `volume_status` - (Optional) The status of the EBS volume.
* `volume_type` - (Optional) The type of the EBS volume.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total_count` - Total number of EBS volumes that satisfy the condition.
* `volumes` - An information list of EBS volumes. Each element contains the following attributes:
  * `availability_zone` - The availability zone in which the EBS volume resides.
  * `create_time` - The time when the EBS volume was created.
  * `instance_id` - The ID of the KEC instance to which the EBS volume is to be attached.
  * `project_id` - The ID of the project.
  * `size` - The capacity of the EBS volume.
  * `volume_category` - The category of the EBS volume.
  * `volume_desc` - The description of the EBS volume.
  * `volume_id` - The ID of the EBS volume.
  * `volume_name` - The name of the EBS volume.
  * `volume_status` - The status of the EBS volume.
  * `volume_type` - The type of the EBS volume.


