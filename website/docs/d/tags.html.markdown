---
subcategory: "Tag"
layout: "ksyun"
page_title: "ksyun: ksyun_tags"
sidebar_current: "docs-ksyun-datasource-tags"
description: |-
  This data source provides a list of tag resources.
---

# ksyun_tags

This data source provides a list of tag resources.

#

## Example Usage

```hcl
data "ksyun_tags" "default" {
  output_file = "output_result"

  # optional
  # eg. key = ["tag_key1", "tag_key2", ...]
  keys = []
  # optional
  # eg. value = ["tag_value1", ...]
  values = []
  # optional
  # eg. resource_type = ["kec-instance", "eip", ...]
  resource_types = []
  # optional
  # eg. key = ["instance_uuid", ...]
  resource_ids = []

}
```

## Argument Reference

The following arguments are supported:

* `keys` - (Optional) A list of tag keys.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_ids` - (Optional) A list of resource ids.
* `resource_types` - (Optional) A list of resource types.
* `values` - (Optional) A list of tag values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tags` - a list of tag.
  * `id` - The ID of the tag.
  * `key` - Tag key.
  * `resource_id` - Resource ID.
  * `resource_type` - Resource type.
  * `value` - Tag value.


