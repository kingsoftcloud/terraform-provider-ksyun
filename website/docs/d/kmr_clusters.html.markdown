---
subcategory: "KMR"
layout: "ksyun"
page_title: "ksyun: ksyun_kmr_clusters"
sidebar_current: "docs-ksyun-datasource-kmr-clusters"
description: |-
  This data source provides a list of KMR clusters.
---

# ksyun_kmr_clusters

This data source provides a list of KMR clusters.

## Example Usage

```hcl
data "ksyun_kmr_clusters" "default" {
  marker = "limit=10&offset=0"
}

output "kmr_clusters" {
  value = data.ksyun_kmr_clusters.default.clusters
}

output "total" {
  value = data.ksyun_kmr_clusters.default.total
}
```

## Argument Reference

* `marker` - (Optional) Pagination marker, e.g., limit=100&offset=0.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference

* `total` - Total number of KMR clusters.
* `clusters` - List of KMR clusters.
  * `cluster_id` - The ID of the cluster.
  * `cluster_name` - The name of the cluster.
  * `main_version` - The main version of the cluster.
  * `enable_eip` - Whether EIP is enabled.
  * `region` - The region of the cluster.
  * `vpc_domain_id` - The VPC domain ID.
  * `charge_type` - The charge type.
  * `cluster_status` - The status of the cluster.
  * `create_time` - The creation time.
  * `update_time` - The update time.
  * `serving_minutes` - The serving minutes.
  * `instance_groups` - List of instance groups.
    * `id` - The ID of the instance group.
    * `instance_group_type` - The type of the instance group.
    * `instance_type` - The instance type.
    * `resource_type` - The resource type.
