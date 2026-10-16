---
subcategory: "KES"
layout: "ksyun"
page_title: "ksyun: ksyun_kes_cluster"
sidebar_current: "docs-ksyun-resource-kes_cluster"
description: |-
  Create a KES cluster (standard / KES7, LaunchCluster v2020-12-15).
---

# ksyun_kes_cluster

Create a KES cluster (standard / KES7, LaunchCluster v2020-12-15).

#

## Example Usage

```hcl
resource "ksyun_kes_cluster" "example" {
  cluster_name      = "my-kes-cluster"
  main_version      = "7.10.0"
  scenario          = "general"
  availability_zone = "cn-beijing-6a"
  charge_type       = "Monthly"
  purchase_time     = 1
  project_id        = 0
  vpc_domain_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  vpc_subnet_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  security_group_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

  instance_groups {
    instance_group_type = "MASTER"
    instance_type       = "ES.I3.2C4G"
    resource_type       = "KEC"
    instance_count      = 3
    volume_type         = "SSD3.0"
    volume_size         = 100
    volume_count        = 1
    system_disk_type    = "SSD3.0"
    system_disk_size    = 50
  }

  instance_groups {
    instance_group_type = "DATA"
    instance_type       = "ES.D4.8C16G"
    resource_type       = "KEC"
    instance_count      = 3
    volume_type         = "SSD3.0"
    volume_size         = 500
    volume_count        = 1
    system_disk_type    = "SSD3.0"
    system_disk_size    = 50
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) The availability zone of the KES cluster.
* `charge_type` - (Required, ForceNew) The charge type of the KES cluster.
* `cluster_name` - (Required, ForceNew) The name of the KES cluster.
* `instance_groups` - (Required, ForceNew) The instance groups of the KES cluster.
* `main_version` - (Required, ForceNew) The main version of the KES cluster, e.g. 7.10.0/7.4.2/6.8.4/5.6.16.
* `project_id` - (Required, ForceNew) The project id of the KES cluster.
* `purchase_time` - (Required, ForceNew) The purchase time of the KES cluster.
* `scenario` - (Required, ForceNew) The scenario of the KES cluster.
* `vpc_domain_id` - (Required, ForceNew) The VPC domain id of the KES cluster.
* `vpc_subnet_id` - (Required, ForceNew) The VPC subnet id of the KES cluster.
* `security_group_id` - (Optional, ForceNew) The security group id of the KES cluster.
* `vpc_epc_subnet_id` - (Optional, ForceNew) The VPC EPC subnet id of the KES cluster.

The `instance_groups` object supports the following:

* `instance_count` - (Required) The instance count of the instance group.
* `instance_group_type` - (Required) The type of the instance group, e.g. MASTER/DATA.
* `instance_type` - (Required) The instance type of the instance group.
* `resource_type` - (Required) The resource type of the instance group.
* `system_disk_size` - (Required) The system disk size (GB) of the instance group.
* `system_disk_type` - (Required) The system disk type of the instance group.
* `volume_count` - (Required) The volume count of the instance group.
* `volume_size` - (Required) The volume size (GB) of the instance group.
* `volume_type` - (Required) The volume type of the instance group.
* `multi_instance_count` - (Optional) The multi instance count of the instance group.
* `raid_type` - (Optional) The raid type of the instance group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cluster_id` - The ID of the KES cluster.
* `cluster_status` - The status of the KES cluster.
* `cluster_type` - The type of the KES cluster.
* `create_time` - The creation time of the KES cluster.
* `enable_eip` - Whether EIP is enabled.
* `proxy_port` - The proxy port of the KES cluster.
* `region` - The region of the KES cluster.
* `serving_minutes` - The serving minutes of the KES cluster.
* `tags` - The tags of the KES cluster.
* `update_time` - The update time of the KES cluster.


## Import

KES cluster can be imported using the id, e.g.

```
$ terraform import ksyun_kes_cluster.example ${cluster_id}
```

