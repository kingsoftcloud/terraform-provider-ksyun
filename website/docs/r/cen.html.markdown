---
subcategory: "CEN"
layout: "ksyun"
page_title: "ksyun: ksyun_cen"
sidebar_current: "docs-ksyun-resource-cen"
description: |-
  Provides a Cen resource.
---

# ksyun_cen

Provides a Cen resource.

#

## Example Usage

```hcl
resource "ksyun_cen" "default" {
  cen_name    = "cen_create"
  description = "zice_create"
}
```

## Argument Reference

The following arguments are supported:

* `cen_name` - (Optional) The name of the cen.
* `description` - (Optional) Cen Description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cen_id` - ID of the cen.
* `create_time` - creation time of the cen.


## Import

Cen can be imported using the `id`, e.g.

```
$ terraform import ksyun_cen.example xxxxxxxx-abc123456
```

