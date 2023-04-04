---
subcategory: "EIP"
layout: "ksyun"
page_title: "ksyun: ksyun_bws_associate"
sidebar_current: "docs-ksyun-resource-bws_associate"
description: |-
  Provides a BWS Association resource for associating EIP with a BWS instance.
---

# ksyun_bws_associate

Provides a BWS Association resource for associating EIP with a BWS instance.

#

## Example Usage

```hcl
resource "ksyun_bws_associate" "default" {
  band_width_share_id = "2af77683-b47e-4634-88ce-fcb95cb65e86"
  allocation_id       = "139134fc-f622-467f-a8b1-c0858dac62ab"
}
```

## Argument Reference

The following arguments are supported:

* `allocation_id` - (Required, ForceNew) ID of the EIP.
* `band_width_share_id` - (Required, ForceNew) ID of the BWS.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `band_width` - bandwidth value.


## Import

# BWS can be imported using the id

```
$ terraform import ksyun_bws_associate.default ${band_width_share_id}:${allocation_id}
```

