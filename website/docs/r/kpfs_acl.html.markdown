---
subcategory: "KPFS"
layout: "ksyun"
page_title: "ksyun: ksyun_kpfs_acl"
sidebar_current: "docs-ksyun-resource-kpfs_acl"
description: |-
  Provides a kpfs acl rule resource.
---

# ksyun_kpfs_acl

Provides a kpfs acl rule resource.

#

## Example Usage

```hcl
resource "ksyun_kpfs_acl" "default" {
  epc_id      = "c6c683f8-5bb4-4747-8516-9a61f01c4bce"
  kpfs_acl_id = "4a42284d7f354e6a9b2d7a3454e0b495"
}
```

## Argument Reference

The following arguments are supported:

* `epc_id` - (Required, ForceNew) The epc instance id.
* `kpfs_acl_id` - (Required, ForceNew) The posix acl rule id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

KPFS ACL rules can be imported using the id, e.g.

```
$ terraform import ksyun_kpfs_acl.example ${epc_id}_${kpfs_acl_id}
```

