---
subcategory: "EIP"
layout: "ksyun"
page_title: "ksyun: ksyun_bws"
sidebar_current: "docs-ksyun-resource-bws"
description: |-
  Provides a BandWidthShare resource.
---

# ksyun_bws

Provides a BandWidthShare resource.

## Example Usage

```hcl
resource "ksyun_bws" "default" {
  line_id     = "5fc2595f-1bfd-481b-bf64-2d08f116d800"
  charge_type = "PostPaidByPeak"
  band_width  = 12
}
```

## Argument Reference

The following arguments are supported:

* `band_width` - (Required) bandwidth value, value range: [1, 15000].
* `charge_type` - (Required, ForceNew) The charge type of the BWS. Valid values: PostPaidByPeak, PostPaidByDay, PostPaidByTransfer, DailyPaidByTransfer.
* `line_id` - (Required, ForceNew) The id of the line.
* `band_width_share_name` - (Optional) name of the BWS.
* `project_id` - (Optional) ID of the project.
* `tags` - (Optional) the tags of the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

BWS can be imported using the id, e.g.

```
$ terraform import ksyun_bws.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

