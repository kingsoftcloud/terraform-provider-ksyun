---
subcategory: "KES8"
layout: "ksyun"
page_title: "ksyun: ksyun_kes_cluster_vector"
sidebar_current: "docs-ksyun-resource-kes_cluster_vector"
description: |-
  Create a KES vector cluster (KES8, 8.x/9.x) via the shared kmr/2021-09-02 LaunchCluster backend switched into KES mode by ClusterType="KES".
---

# ksyun_kes_cluster_vector

Create a KES vector cluster (KES8, 8.x/9.x) via the shared kmr/2021-09-02 LaunchCluster backend switched into KES mode by ClusterType="KES".

KES8 clusters run Elasticsearch 8.x / 9.x and require Xpack credentials. The
authoritative creation path (confirmed by a Ksyun console capture + a production
probe) is the SHARED kmr/2021-09-02 LaunchCluster action, with the body field
ClusterType="KES" switching the backend from KMR mode into KES8 mode. The SDK
ships no typed client for that version, so this resource builds a custom request
struct and dispatches it through the KES client's generic signing channel (see
service_ksyun_kes_vector.go).

On the kmr/2021-09-02 backend, placement (AvailabilityZone, VpcSubnetId, VpcId)
and the per-group index/ChargeType live INSIDE each InstanceGroups entry, and
VpcSubnetId is also sent at top level. PurchaseTime and VpcEpcSubnetId are not
part of the KES8 shape and are therefore not modelled here.

#

## Example Usage

```hcl
resource "ksyun_kes_cluster_vector" "example" {
  cluster_name      = "my-kes-vector-cluster"
  main_version      = "8.19.13"
  distribution      = "ElasticSearch"
  scenario          = "General"
  charge_type       = "Monthly"
  project_id        = 0
  vpc_domain_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  vpc_subnet_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  security_group_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

  # Required for any 8.x / 9.x cluster.
  xpack_username = "elastic"
  xpack_password = "your-xpack-password"

  # Optional KES8 switches (defaults shown).
  xpack_enabled             = true
  enable_best_delivery_mode = false
  enable_eip                = false
  instance_login_type       = "later"
  services                  = ["elasticsearch"]
  plugins                   = ["repository-ks3"]

  instance_groups {
    instance_group_type     = "MASTER"
    instance_type           = "ES.I3.2C4G"
    instance_type_code      = "ES.I3.2C4G"
    resource_type           = "KEC"
    instance_count          = 3
    instance_group_index    = 1
    availability_zone       = "cn-beijing-6a"
    availability_zone_index = 1
    vpc_id                  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    vpc_subnet_id           = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    charge_type             = "Monthly"
    volume_type             = "SSD3.0"
    volume_size             = 100
    volume_count            = 1
    system_disk_type        = "SSD3.0"
    system_disk_size        = 50
  }

  instance_groups {
    instance_group_type     = "DATA"
    instance_type           = "ES.D4.8C16G"
    instance_type_code      = "ES.D4.8C16G"
    resource_type           = "KEC"
    instance_count          = 3
    instance_group_index    = 2
    availability_zone       = "cn-beijing-6a"
    availability_zone_index = 1
    vpc_id                  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    vpc_subnet_id           = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    charge_type             = "Monthly"
    volume_type             = "SSD3.0"
    volume_size             = 500
    volume_count            = 1
    system_disk_type        = "SSD3.0"
    system_disk_size        = 50
  }
}
```

## Argument Reference

The following arguments are supported:

* `charge_type` - (Required, ForceNew) The charge type of the KES vector cluster, e.g. Monthly/Minutely.
* `cluster_name` - (Required, ForceNew) The name of the KES vector cluster.
* `instance_groups` - (Required, ForceNew) The instance groups of the KES vector cluster. Placement (availability_zone/vpc_subnet_id/vpc_id) and the per-group index/charge_type are set per group.
* `main_version` - (Required, ForceNew) The main version of the KES vector cluster, e.g. 8.19.13/9.x. Xpack credentials are required for any 8.x or 9.x cluster.
* `project_id` - (Required, ForceNew) The project id of the KES vector cluster.
* `scenario` - (Required, ForceNew) The scenario of the KES vector cluster, e.g. General.
* `vpc_domain_id` - (Required, ForceNew) The VPC domain id of the KES vector cluster.
* `vpc_subnet_id` - (Required, ForceNew) The VPC subnet id of the KES vector cluster. Also sent per instance group on the kmr/2021-09-02 backend.
* `xpack_password` - (Required, ForceNew) The Xpack password (at least 6 characters). Required for KES8 (8.x/9.x) clusters.
* `xpack_username` - (Required, ForceNew) The Xpack username. Required for KES8 (8.x/9.x) clusters.
* `distribution` - (Optional, ForceNew) The distribution of the KES vector cluster. Defaults to ElasticSearch.
* `enable_best_delivery_mode` - (Optional, ForceNew) Whether best delivery mode is enabled.
* `enable_eip` - (Optional, ForceNew) Whether EIP is enabled for the cluster.
* `instance_login_type` - (Optional, ForceNew) The instance login type, e.g. 'later' (no password) or 'password'.
* `instance_password` - (Optional, ForceNew) The instance login password. Only sent when instance_login_type is 'password'.
* `plugins` - (Optional, ForceNew) The plugins to install, e.g. repository-ks3. Only sent when set.
* `security_group_id` - (Optional, ForceNew) The security group id of the KES vector cluster.
* `services` - (Optional, ForceNew) The services to deploy. Defaults to ["elasticsearch"] when omitted.
* `xpack_enabled` - (Optional, ForceNew) Whether Xpack is enabled. Defaults to true for KES8 clusters.

The `instance_groups` object supports the following:

* `availability_zone` - (Required) The availability zone of the instance group, e.g. cn-beijing-6a.
* `instance_count` - (Required) The instance count of the instance group.
* `instance_group_type` - (Required) The type of the instance group, e.g. MASTER/DATA.
* `instance_type` - (Required) The instance type of the instance group.
* `resource_type` - (Required) The resource type of the instance group, e.g. KEC.
* `system_disk_size` - (Required) The system disk size (GB) of the instance group.
* `system_disk_type` - (Required) The system disk type of the instance group, e.g. SSD3.0.
* `volume_count` - (Required) The volume count of the instance group.
* `volume_size` - (Required) The volume size (GB) of the instance group.
* `volume_type` - (Required) The volume type of the instance group, e.g. SSD3.0.
* `availability_zone_index` - (Optional) The availability zone index of the instance group.
* `charge_type` - (Optional) The charge type of the instance group, e.g. Monthly/Minutely.
* `instance_group_index` - (Optional) The index of the instance group.
* `instance_type_code` - (Optional) The instance type code of the instance group, e.g. ES.X6.4C8G.
* `multi_instance_count` - (Optional) The multi instance count of the instance group.
* `raid_type` - (Optional) The raid type of the instance group.
* `vpc_id` - (Optional) The VPC id of the instance group.
* `vpc_subnet_id` - (Optional) The VPC subnet id of the instance group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cluster_id` - The ID of the KES vector cluster.
* `cluster_status` - The status of the KES vector cluster.
* `cluster_type` - The type of the KES vector cluster (KES).
* `create_time` - The creation time of the KES vector cluster.
* `proxy_port` - The proxy port of the KES vector cluster.
* `region` - The region of the KES vector cluster.
* `serving_minutes` - The serving minutes of the KES vector cluster.
* `tags` - The tags of the KES vector cluster.
* `update_time` - The update time of the KES vector cluster.


## Import

KES vector cluster can be imported using the id, e.g.

```
$ terraform import ksyun_kes_cluster_vector.example ${cluster_id}
```

