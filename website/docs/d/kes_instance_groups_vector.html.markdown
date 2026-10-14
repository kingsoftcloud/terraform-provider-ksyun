---
subcategory: "KES8"
layout: "ksyun"
page_title: "ksyun: ksyun_kes_instance_groups_vector"
sidebar_current: "docs-ksyun-datasource-kes_instance_groups_vector"
description: |-
  Provides a KES vector instance groups (KES8) data source.
---

# ksyun_kes_instance_groups_vector

Provides a KES vector instance groups (KES8) data source.

#

## Example Usage

```hcl
provider "ksyun" {
  region = "cn-beijing-6"
}

data "ksyun_kes_instance_groups_vector" "default" {
  cluster_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

output "first_instance_group_id" {
  value = data.ksyun_kes_instance_groups_vector.default.instance_groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the KES cluster whose instance groups are queried.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_groups` - List of instance groups of the cluster.
  * `avalability_zone` - The availability zone.
  * `id` - The ID of the instance group.
  * `instance_count` - The instance count.
  * `instance_group_type` - The type of the instance group.
  * `instance_type_code` - The instance type code.
  * `instances` - List of instances.
    * `id` - The ID of the instance.
    * `instance_group_id` - The instance group ID of the instance.
    * `instance_id` - The instance ID.
    * `instance_name` - The instance name.
    * `internal_ip` - The internal IP.
    * `volumes` - List of volume identifiers.
  * `multi_instance_count` - The multi instance count.
  * `resource_type` - The resource type.
  * `volume_count` - The volume count.
  * `volume_size` - The volume size.
  * `volume_type` - The volume type.
  * `vpc_id` - The VPC ID.
  * `vpc_subnet_id` - The VPC subnet ID.


