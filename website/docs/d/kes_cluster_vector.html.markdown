---
subcategory: "KES8"
layout: "ksyun"
page_title: "ksyun: ksyun_kes_cluster_vector"
sidebar_current: "docs-ksyun-datasource-kes_cluster_vector"
description: |-
  Provides a KES vector cluster (KES8) data source.
---

# ksyun_kes_cluster_vector

Provides a KES vector cluster (KES8) data source.

#

## Example Usage

```hcl
provider "ksyun" {
  region = "cn-beijing-6"
}

data "ksyun_kes_cluster_vector" "default" {
  cluster_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

output "cluster_name" {
  value = data.ksyun_kes_cluster_vector.default.clusters.0.cluster_name
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the KES cluster to query.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `clusters` - A list containing a single KES cluster detail.
  * `availability_zone` - The availability zone of the cluster.
  * `charge_type` - The charge type.
  * `cluster_id` - The ID of the cluster.
  * `cluster_name` - The name of the cluster.
  * `cluster_status` - The status of the cluster.
  * `cluster_type` - The type of the cluster.
  * `create_time` - The creation time.
  * `enable_eip` - Whether EIP is enabled.
  * `instance_groups` - List of instance groups.
    * `id` - The ID of the instance group.
    * `instance_group_type` - The type of the instance group.
    * `instance_type` - The instance type.
    * `resource_type` - The resource type.
    * `volume_size` - The volume size.
    * `volume_type` - The volume type.
  * `main_version` - The main version of the cluster.
  * `proxy_port` - The proxy port of the cluster.
  * `region` - The region of the cluster.
  * `serving_minutes` - The serving minutes.
  * `tags` - List of tags.
    * `tag_key` - The key of the tag.
    * `tag_value` - The value of the tag.
  * `update_time` - The update time.
  * `vpc_domain_id` - The VPC domain ID.
  * `vpc_subnet_id` - The VPC subnet ID.


