---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_direct_connects"
sidebar_current: "docs-ksyun-datasource-direct_connects"
description: |-
  This data source provides a list of Direct Connect resources.
---

# ksyun_direct_connects

This data source provides a list of Direct Connect resources.

## Example Usage

```hcl
data "ksyun_direct_connects" "test" {
  ids        = []
  name_regex = ".*test.*"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Direct Connect IDs.
* `name_regex` - (Optional) A regex string to filter results by Direct Connect name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `direct_connects` - It is a nested type which documented below.
  * `band_width` - Band Width.
  * `create_time` - creation time of the Direct Connect.
  * `customer_location` - Customer location of the Direct Connect.
  * `direct_connect_id` - ID of the Direct Connect.
  * `direct_connect_name` - name of the Direct Connect.
  * `distance` - Distance.
  * `id` - ID of the Direct Connect.
  * `name` - name of the Direct Connect.
  * `pop_location` - Pop Location.
  * `state` - State.
  * `type` - Type.
  * `vlan` - Vlan.
  * `vpc_noc_id` - Vpc Noc ID.
* `total_count` - Total number of Direct Connect that satisfy the condition.


