---
subcategory: "KPFS"
layout: "ksyun"
page_title: "ksyun: ksyun_kpfs_client_install"
sidebar_current: "docs-ksyun-datasource-kpfs_client_install"
description: |-
  Query KPFS client installation package information and mount IP by file system ID.
---

# ksyun_kpfs_client_install

Query KPFS client installation package information and mount IP by file system ID.

#

## Example Usage

```hcl
data "ksyun_kpfs_client_install" "default" {
  output_file = "output_result"
  id          = "b7449ea1d57f428595f7c68a1fbeeafd"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) File system ID.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - A list of KPFS clusters.
  * `cluster_data_ip` - The data IP of the KPFS cluster.
  * `download_url` - The download URL for the KPFS client.
  * `kernel_version` - The kernel version supported by the KPFS client.
  * `nic_driver` - The NIC driver version supported by the KPFS client.
  * `os_version` - The OS version supported by the KPFS client.


