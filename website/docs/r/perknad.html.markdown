---
subcategory: "KNAD"
layout: "ksyun"
page_title: "ksyun: ksyun_perknad"
sidebar_current: "docs-ksyun-resource-perknad"
description: |-
  Provides a PerPay KNAD (PerKnad) resource.
---

# ksyun_perknad

Provides a PerPay KNAD (PerKnad) resource.

#

## Example Usage

```hcl
# Create a perpay knad

resource "ksyun_perknad" "default" {
  ip_count   = 10
  max_band   = 30
  knad_name  = "ksc_kad"
  project_id = "0"
}
```

## Argument Reference

The following arguments are supported:

* `ip_count` - (Required) the max ip count that can bind to the PerKnad, value range: [10, 100].
* `max_band` - (Required) the max protection band of the PerKnad.
* `knad_name` - (Optional) the name of the PerKnad.
* `project_id` - (Optional) The id of the project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `knad_id` - the ID of the PerKnad.


## Import

PerKnad can be imported using the id, e.g.

```
$ terraform import ksyun_perknad.default knad67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

