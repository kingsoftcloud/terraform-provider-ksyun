---
subcategory: "Redis"
layout: "ksyun"
page_title: "ksyun: ksyun_redis_instances"
sidebar_current: "docs-ksyun-datasource-redis_instances"
description: |-
  Provides a list of Redis resources in the current region.
---

# ksyun_redis_instances

Provides a list of Redis resources in the current region.

#

## Example Usage

```hcl
data "ksyun_redis_instances" "default" {
  output_file    = "output_result1"
  fuzzy_search   = ""
  iam_project_id = ""
  cache_id       = ""
  vnet_id        = ""
  vpc_id         = ""
  name           = ""
  vip            = ""
}
```

## Argument Reference

The following arguments are supported:

* `cache_id` - (Optional) The ID of the instance.
* `fuzzy_search` - (Optional) fuzzy filter by name / VIP / ID.
* `iam_project_id` - (Optional) The project instance belongs to.
* `name` - (Optional) The name of instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `vip` - (Optional) Private IP address of the instance.
* `vnet_id` - (Optional) The ID of subnet.
* `vpc_id` - (Optional) The ID of VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - An information list of instances. Each element contains the following attributes:
  * `az` - Availability zone.
  * `bill_type` - Bill type.
  * `cache_id` - The ID of the instance.
  * `create_time` - creation time of instance.
  * `eip_ro` - EIP address of read-only node.
  * `eip` - EIP address.
  * `engine` - Engine.
  * `iam_project_id` - project id.
  * `iam_project_name` - project name.
  * `mode` - The KVStore instance system architecture.
  * `name` - The name of instance.
  * `net_type` - Type of network.
  * `order_type` - Order type.
  * `parameters` - parameters of instance.
  * `port` - port number.
  * `protocol` - protocol of instance.
  * `readonly_node` - A list of read-only nodes.
    * `create_time` - Creation time.
    * `instance_id` - The ID of instance.
    * `ip` - Private IP.
    * `name` - The name of instance.
    * `port` - Port number.
    * `proxy` - Role of node.
    * `status` - Status.
  * `region` - Region.
  * `service_begin_time` - service begin time.
  * `service_end_time` - service end time.
  * `service_status` - service status.
  * `shard_num` - Shard num.
  * `shard_size` - Shard memory size.
  * `size` - The size of instance.
  * `source` - size of source.
  * `status` - Status.
  * `timezone` - Auto backup time zone.
  * `timing_switch` - Switch auto backup.
  * `vip` - Private IP address of the instance.
  * `vnet_id` - The id of subnet linked to the instance.
  * `vpc_id` - The id of VPC linked to the instance.
* `total_count` - Total number of Redis instances that satisfy the condition.


