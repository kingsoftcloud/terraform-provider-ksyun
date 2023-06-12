---
subcategory: "KCE"
layout: "ksyun"
page_title: "ksyun: ksyun_kce_instance_images"
sidebar_current: "docs-ksyun-datasource-kce_instance_images"
description: |-
  This data source providers a list of available instance image which support kce.
---

# ksyun_kce_instance_images

This data source providers a list of available instance image which support kce.

#

## Example Usage

```hcl
data "ksyun_kce_instance_images" "default" {
  output_file = "output_result"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `image_set` - a list of images.
  * `image_id` - The ID of the image.
  * `image_name` - The name of the image.
  * `image_type` - The type of the image.


