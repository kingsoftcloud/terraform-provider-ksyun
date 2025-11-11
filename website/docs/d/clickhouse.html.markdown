---
subcategory: "Clickhouse"
layout: "ksyun"
page_title: "ksyun: ksyun_clickhouse"
sidebar_current: "docs-ksyun-datasource-clickhouse"
description: |-
  Query Clickhouse instance information
---

# ksyun_clickhouse

Query Clickhouse instance information

#

## Example Usage

```hcl
data "ksyun_clickhouse" "default" {

  instance_id = "instance_id"

  product_type = "product_type"
  project_ids  = "project_ids"
  tag_id       = "tag_id"

  fuzzy_search = "fuzzy_search"

  offset = 0
  limit  = 10
}
```

## Argument Reference

The following arguments are supported:

* `fuzzy_search` - (Optional) Fuzzy search filter that matches instance name, VIP, or instance ID.
* `instance_id` - (Optional) The ClickHouse instance ID. When provided, returns detailed information for that specific instance; otherwise returns a list of all instances.
* `limit` - (Optional) The maximum number of records to return per page. Default is 10.
* `offset` - (Optional) The starting offset for pagination. Default is 0 (first page).
* `output_file` - (Optional) File name where to save data source results (after running terraform plan).
* `product_type` - (Optional) The product type of the instance. Valid values: 'ClickHouse_Single' (single replica) or 'ClickHouse' (high availability).
* `project_ids` - (Optional) Comma-separated list of project IDs to filter instances.
* `tag_id` - (Optional) Filter instances by tag ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - List of ClickHouse instances.
  * `admin_user` - Admin user name.
  * `area` - Instance area configuration.
    * `master` - Master area.
    * `standby` - Standby area.
  * `az` - Availability zone.
  * `bill_type` - The billing type of the instance.
  * `cpu` - CPU cores.
  * `create_date` - Creation time.
  * `direct_connection_vips` - Direct connection VIPs.
  * `ebs_size` - Disk size in GB.
  * `ebs_type` - Disk type.
  * `engine_version` - Engine version.
  * `engine` - Database engine.
  * `hot_and_cold` - Hot and cold configuration.
  * `http_port` - HTTP port.
  * `instance_config` - Instance configuration.
  * `instance_id` - Instance ID.
  * `instance_name` - Instance name.
  * `mem` - Memory size in GB.
  * `multiaz` - Multi-Availability Zone deployment configuration flag.
  * `network_type` - The network type of the instance. Currently supports 'VPC'.
  * `node_num` - Number of nodes.
  * `product_id` - Product ID.
  * `product_type_name` - Product type name.
  * `product_type` - Product type.
  * `product_what` - The product category identifier.
  * `project_id` - Project ID.
  * `project_name` - Project name.
  * `region` - Region.
  * `replicas` - Number of replicas.
  * `security_group_desc` - Security group description.
  * `security_group_id` - Security group ID.
  * `security_group_name` - Security group name.
  * `service_end_time` - Service end time.
  * `shard_list` - Instance shard list.
    * `id` - Shard ID.
    * `name` - Shard name.
  * `status_name` - The display name of the instance status (Chinese).
  * `status` - Instance status.
  * `storage_size` - Storage size.
  * `storage_type` - Storage type.
  * `subnet_id` - Subnet ID.
  * `tags` - Instance tags.
    * `tag_id` - Tag ID.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
  * `tcp_port` - TCP port.
  * `update_date` - Update time.
  * `used_storage_size` - Used storage size.
  * `user_id` - User ID.
  * `vip` - Virtual IP address.
  * `vpc_id` - VPC ID.
* `total_count` - Total number of resources that satisfy the condition.


