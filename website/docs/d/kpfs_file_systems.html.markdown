---
subcategory: "KPFS"
layout: "ksyun"
page_title: "ksyun: ksyun_kpfs_file_systems"
sidebar_current: "docs-ksyun-datasource-kpfs_file_systems"
description: |-
  This data source provides a list of kpfs fileSystem resources according to their fileSystem ID, name.
---

# ksyun_kpfs_file_systems

This data source provides a list of kpfs fileSystem resources according to their fileSystem ID, name.

#

## Example Usage

```hcl
data "ksyun_kpfs_file_systems" "default" {
  output_file = "output_result"
  id          = "fileSystemId"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) File system ID.
* `output_file` - (Optional) File name where to save data source results.
* `page_num` - (Optional) Page number for pagination.
* `page_size` - (Optional) Page size for pagination.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - It is a nested type which documented below.
  * `status` - Current state: creating/using/upgrading/renewing/shutdown.
* `total_count` - The total number of matching file systems.


