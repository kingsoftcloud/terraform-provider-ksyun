---
subcategory: "KNAD"
layout: "ksyun"
page_title: "ksyun: ksyun_knad_associate"
sidebar_current: "docs-ksyun-resource-knad_associate"
description: |-
  Provides a Knad Association resource for associating EIP with a KNAD instance.
---

# ksyun_knad_associate

Provides a Knad Association resource for associating EIP with a KNAD instance.

#

## Example Usage

```hcl
resource "ksyun_knad_associate" "default" {
  knad_id = "knadba4d704f-35b1-3354-8d0f-65adc04796b9"
  ip      = ["88.88.88.46", "88.88.88.51"]
}
```

## Argument Reference

The following arguments are supported:

* `ip` - (Optional) the binding ips.
* `knad_id` - (Optional, ForceNew) the ID of the Knad.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

# Knad Association can be imported using the id

```
$ terraform import ksyun_knad_associate.default ${knad_id}
```

