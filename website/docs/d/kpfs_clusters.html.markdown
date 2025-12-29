---
subcategory: "KPFS"
layout: "ksyun"
page_title: "ksyun: ksyun_kpfs_clusters"
sidebar_current: "docs-ksyun-datasource-kpfs_clusters"
description: |-
  Query available storage cluster information by UID.
---

# ksyun_kpfs_clusters

Query available storage cluster information by UID.

#

## Example Usage

```hcl
data "ksyun_kpfs_clusters" "default" {
  output_file     = "output_result"
  region          = "cn-qingyangtest-1"
  s_roce_cluster  = "QYYC01-Sroce-Cluster-01"
  store_class     = "KPFS-P-S01"
  store_pool_type = "KPFS-P1"
  avail_zone      = "cn-qingyangtest-1a"
}
```

## Argument Reference

The following arguments are supported:

* `avail_zone` - (Optional) The availability zone of the KPFS cluster.
* `cluster_code` - (Optional) The unique code of the KPFS cluster.
* `output_file` - (Optional) File name where to save data source results.
* `region` - (Optional) The region of the KPFS cluster.
* `s_roce_cluster` - (Optional) The SRoCE cluster name of the KPFS cluster.
* `store_class` - (Optional) The storage classes supported by the KPFS cluster.
* `store_pool_type` - (Optional) The storage pool type of the KPFS cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - A list of KPFS clusters.
  * `avail_zone` - The availability zone of the KPFS cluster.
  * `cluster_code` - The unique code of the KPFS cluster.
  * `region` - The region of the KPFS cluster.
  * `s_roce_cluster` - The SRoCE cluster name of the KPFS cluster.
  * `store_classes` - The storage classes supported by the KPFS cluster.
  * `store_pool_type` - The storage pool type of the KPFS cluster.


