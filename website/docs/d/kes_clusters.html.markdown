---
subcategory: "KES"
layout: "ksyun"
page_title: "ksyun: ksyun_kes_clusters"
sidebar_current: "docs-ksyun-datasource-kes_clusters"
description: |-
  Provides a KES Clusters data source.
---

# ksyun_kes_clusters

Provides a KES Clusters data source.

#

## Example Usage

```hcl
provider "ksyun" {
  region = "cn-beijing-6"
}

data "ksyun_kes_clusters" "default" {
  marker = "limit=10&offset=0"
}

output "total" {
  value = data.ksyun_kes_clusters.default.total
}
```

## Argument Reference

The following arguments are supported:

* `marker` - (Optional) Pagination marker, e.g., limit=100&offset=0.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `clusters` - List of KES clusters.
  * `availability_zone` - The availability zone of the cluster.
  * `charge_type` - The charge type.
  * `cluster_id` - The ID of the cluster.
  * `cluster_name` - The name of the cluster.
  * `cluster_status` - The status of the cluster.
  * `create_time` - The creation time.
  * `enable_eip` - Whether EIP is enabled.
  * `instance_groups` - List of instance groups.
    * `id` - The ID of the instance group.
    * `instance_group_type` - The type of the instance group.
    * `instance_type` - The instance type.
    * `resource_type` - The resource type.
  * `main_version` - The main version of the cluster.
  * `region` - The region of the cluster.
  * `serving_minutes` - The serving minutes.
  * `update_time` - The update time.
  * `vpc_domain_id` - The VPC domain ID.
  * `vpc_subnet_id` - The VPC subnet ID.
* `total` - Total number of KES clusters.


