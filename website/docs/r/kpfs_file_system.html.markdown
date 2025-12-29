---
subcategory: "KPFS"
layout: "ksyun"
page_title: "ksyun: ksyun_kpfs_file_system"
sidebar_current: "docs-ksyun-resource-kpfs_file_system"
description: |-
  Create a kpfs file system.
---

# ksyun_kpfs_file_system

Create a kpfs file system.

#

## Example Usage

```hcl
resource "ksyun_kpfs_file_system" "example" {
  file_system_name = "examplefs9"
  region           = "cn-qingyangtest-1"
  avail_zone       = "cn-qingyangtest-1a"
  charge_type      = "dailySettlement"
  project_id       = "0"
  purchase_time    = "-1"
  store_class      = "KPFS-P-S01"
  capacity         = 102
  chunk_size_type  = "Balanced"
  cluster_code     = "7c2cd53b-8eec-440d-a450-18b7db2b0040-5"
}
```

## Argument Reference

The following arguments are supported:

* `avail_zone` - (Required) The availability zone of the file system.
* `capacity` - (Required) The capacity of the file system.
* `charge_type` - (Required) The charge type of the file system.
* `chunk_size_type` - (Required) The chunk size type of the file system.
* `cluster_code` - (Required) The cluster code of the file system.
* `file_system_name` - (Required, ForceNew) The name of the file system.
* `project_id` - (Required) The project id of the file system.
* `purchase_time` - (Required) The purchase time of the file system.
* `region` - (Required) The region of the file system.
* `store_class` - (Required) The store class of the file system.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

KPFS filesystem can be imported using the id, e.g.

```
$ terraform import ksyun_kpfs_file_system.example ${filesystem_id}
```

