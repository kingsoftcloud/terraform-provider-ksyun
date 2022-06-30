---
layout: "ksyun"
page_title: "Ksyun: ksyun_tags"
sidebar_current: "docs-ksyun-datasource-tags"
description: |-
  Provides a list of tag resources in the current region.
---

# ksyun_tags

This data source provides a list of tag resources.

## Example Usage

```hcl
data "ksyun_tags" "default" {
  output_file="output_result"

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

* `keys` - (Optional) A list of tag keys
* `values` - (Optional) A list of tag values
* `resource_types` - (Optional) A list of resource types
* `resource_ids` - (Optional) A list of resource ids
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
