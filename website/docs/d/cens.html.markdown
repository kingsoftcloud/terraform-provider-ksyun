---
subcategory: "CEN"
layout: "ksyun"
page_title: "ksyun: ksyun_cens"
sidebar_current: "docs-ksyun-datasource-cens"
description: |-
  This data source provides a list of Subnet resources according to their Subnet ID, name and the VPC they belong to.
---

# ksyun_cens

This data source provides a list of Subnet resources according to their Subnet ID, name and the VPC they belong to.

#

## Example Usage

```hcl
data "ksyun_cens" "default" {
  output_file = "output_result"
  ids         = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Cen IDs, all the Cens belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter results by cen name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cens` - An information list of cens. Each element contains the following attributes:
  * `cen_id` - ID of the cen.
  * `cen_name` - Name of the cen.
  * `create_time` - creation time of the cen.
  * `description` - The description of cen.
  * `id` - ID of the cen.
* `total_count` - Total number of cens that satisfy the condition.


