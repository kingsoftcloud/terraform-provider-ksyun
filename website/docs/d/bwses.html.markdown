---
subcategory: "EIP"
layout: "ksyun"
page_title: "ksyun: ksyun_bwses"
sidebar_current: "docs-ksyun-datasource-bwses"
description: |-
  This data source provides a list of BWS resources (BandWidthShare) according to their BWS ID.
---

# ksyun_bwses

This data source provides a list of BWS resources (BandWidthShare) according to their BWS ID.

## Example Usage

```hcl
data "ksyun_bwses" "default" {
  output_file = "output_result"
  ids         = ["c7b2ba05-9302-4933-8588-a66f920ff57d"]
}
```

## Argument Reference

The following arguments are supported:

* `allocation_ids` - (Optional) One or more ids of the EIPs in the BWS.
* `ids` - (Optional) A list of BWS IDs, all the BWSs belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter results by BWS name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `project_ids` - (Optional) One or more project IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `band_width_shares` - It is a nested type which documented below.
  * `allocation_ids` - a list of EIP IDs which associated to BWS.
  * `band_width_share_id` - ID of the BWS.
  * `band_width_share_name` - name of the BWS.
  * `band_width` - bandwidth value.
  * `create_time` - creation time of the BWS.
  * `id` - ID of the BWS.
  * `line_id` - ID of the BWS line.
  * `name` - name of the BWS.
  * `project_id` - ID of project.
* `total_count` - Total number of BWS that satisfy the condition.


