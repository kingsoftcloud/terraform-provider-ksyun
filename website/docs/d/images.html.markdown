---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_images"
sidebar_current: "docs-ksyun-datasource-images"
description: |-
  This data source providers a list of available image resources according to their availability zone, image ID and other fields.
---

# ksyun_images

This data source providers a list of available image resources according to their availability zone, image ID and other fields.

#

## Example Usage

```hcl
data "ksyun_images" "default" {
  output_file  = "output_result"
  is_public    = true
  image_source = "system"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of image IDs.
* `image_source` - (Optional) Valid values are import, copy, share, extend, system.
* `is_public` - (Optional) If ksyun provide the image.
* `name_regex` - (Optional) A regex string to filter resulting images by name. (Such as: `^CentOS 7.[1-2] 64` means CentOS 7.1 of 64-bit operating system or CentOS 7.2 of 64-bit operating system, "^Ubuntu 16.04 64" means Ubuntu 16.04 of 64-bit operating system).
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `platform` - (Optional) Platform type of the image system.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `images` - It is a nested type which documented below.
  * `cloud_init_support` - Whether support cloud-init.
  * `creation_date` - Time of creation.
  * `image_id` - The ID of image.
  * `image_source` - Image source of the image.
  * `image_state` - Status of the image.
  * `instance_id` - the id of the instance which the image based on.
  * `ipv6_support` - Whether support ipv6.
  * `is_cloud_market` - Whether image is from cloud market or not.
  * `is_modify_type` - Whether support live upgrade.
  * `is_npe` - whether networking enhancement is support or not.
  * `is_public` - If ksyun provide the image.
  * `name` - Display name of the image.
  * `platform` - Platform type of the image system.
  * `progress` - image creation progress percentage.
  * `real_image_id` - The real id of the image.
  * `sys_disk` - size of system disk.
  * `user_category` - User defined category.
* `total_count` - Total number of image that satisfy the condition.


