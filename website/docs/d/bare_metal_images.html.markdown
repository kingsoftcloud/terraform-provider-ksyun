---
subcategory: "Bare Metal"
layout: "ksyun"
page_title: "ksyun: ksyun_bare_metal_images"
sidebar_current: "docs-ksyun-datasource-bare_metal_images"
description: |-
  This data source provides a list of Bare Metal Image resources according to their Bare Metal Image ID.
---

# ksyun_bare_metal_images

This data source provides a list of Bare Metal Image resources according to their Bare Metal Image ID.

#

## Example Usage

```hcl
# Get  bare metal_images

data "ksyun_bare_metal_images" "default" {
  output_file = "output_result"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Bare Metal Images IDs, all the Bare Metal Images belong to this region will be retrieved if the ID is `""`.
* `image_type` - (Optional) A list of Bare Metal Images Types.
* `name_regex` - (Optional) A regex string to filter results by name of Bare Metal Image.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `images` - It is a nested type which documented below.
  * `create_time` - The time of creation for Bera Metal Image.
  * `image_id` - ID of the Image.
  * `image_name` - name of the Image.
  * `image_type` - type of the Image.
  * `os_name` - OS name of the Image.
  * `os_type` - OS type of the Image.
* `total_count` - Total number of Bare Metal Images that satisfy the condition.


