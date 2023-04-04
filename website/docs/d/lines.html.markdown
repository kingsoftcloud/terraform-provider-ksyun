---
subcategory: "EIP"
layout: "ksyun"
page_title: "ksyun: ksyun_lines"
sidebar_current: "docs-ksyun-datasource-lines"
description: |-
  This data source provides a list of line resources supported.
---

# ksyun_lines

This data source provides a list of line resources supported.

#

## Example Usage

```hcl
data "ksyun_lines" "default" {
  output_file = "output_result"
  line_name   = "BGP"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of lines, all the lines belong to this region will be retrieved if the ID is `""`.
* `line_name` - (Optional) Name of the line.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `lines` - All the lines according the argument.
  * `line_id` - ID of the line.
  * `line_name` - Name of the line.
  * `line_type` - Type of the line.
* `total_count` - Total number of lines that satisfy the condition.


