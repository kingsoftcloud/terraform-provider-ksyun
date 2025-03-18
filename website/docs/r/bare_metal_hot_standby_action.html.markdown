---
subcategory: "Bare Metal"
layout: "ksyun"
page_title: "ksyun: ksyun_bare_metal_hot_standby_action"
sidebar_current: "docs-ksyun-resource-bare_metal_hot_standby_action"
description: |-
  Provides bare metal use hot standby action.
---

# ksyun_bare_metal_hot_standby_action

Provides bare metal use hot standby action.

#

## Example Usage

```hcl
resource "ksyun_bare_metal_hot_standby_action" "foo" {
  host_id = "epc_id"
  hot_standby {
    hot_stand_by_host_id = "hot_standby_id"
    retain_instance_info = "Notretain"
  }
}
```

## Argument Reference

The following arguments are supported:

* `host_id` - (Required, ForceNew) The id of epc for hot standby.
* `hot_standby` - (Required, ForceNew) Indicate the hot standby to instead the master Host.

The `hot_standby` object supports the following:

* `hot_stand_by_host_id` - (Required, ForceNew) The id of hot standby.
* `retain_instance_info` - (Optional, ForceNew) Whether retain the instance info. Valid Values: `RetainPrivateIP` `Notretain`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Dont Allow to Import

