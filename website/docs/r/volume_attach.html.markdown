---
subcategory: "Volume(EBS)"
layout: "ksyun"
page_title: "ksyun: ksyun_volume_attach"
sidebar_current: "docs-ksyun-resource-volume_attach"
description: |-
  Provides an EBS attachment resource.
---

# ksyun_volume_attach

Provides an EBS attachment resource.

#

## Example Usage

```hcl
resource "ksyun_volume_attach" "default" {
  volume_id            = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  instance_id          = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  delete_with_instance = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the KEC instance to which the EBS volume is to be attached.
* `volume_id` - (Required, ForceNew) The ID of the EBS volume.
* `delete_with_instance` - (Optional) Specifies whether to delete the EBS volume when the KEC instance to which it is attached is deleted. Default value: false.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `availability_zone` - The availability zone in which the EBS volume resides.
* `create_time` - The time when the EBS volume was created.
* `project_id` - The ID of the project.
* `size` - The capacity of the EBS volume, in GB.
* `volume_category` - The category of the EBS volume.
* `volume_desc` - The description of the EBS volume.
* `volume_name` - The name of the EBS volume.
* `volume_status` - The status of the EBS volume.
* `volume_type` - The type of the EBS volume.


## Import

EBS volume can be imported using the `id`, e.g.

```
$ terraform import ksyun_volume.default $volume_id:$instance_id
```

