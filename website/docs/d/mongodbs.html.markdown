---
subcategory: "MongoDB"
layout: "ksyun"
page_title: "ksyun: ksyun_mongodbs"
sidebar_current: "docs-ksyun-datasource-mongodbs"
description: |-
  This data source provides a list of MongoDB resources according to their name, Instance ID, Subnet ID, VPC ID and the Project ID they belong to .
---

# ksyun_mongodbs

This data source provides a list of MongoDB resources according to their name, Instance ID, Subnet ID, VPC ID and the Project ID they belong to .

#

## Example Usage

```hcl
data "ksyun_mongodbs" "default" {
  output_file    = "output_result"
  iam_project_id = ""
  instance_id    = ""
  vnet_id        = ""
  vpc_id         = ""
  name           = ""
  vip            = ""
}
```

## Argument Reference

The following arguments are supported:

* `iam_project_id` - (Optional) The project instance belongs to.
* `instance_id` - (Optional) The id of MongoDB, all the MongoDBs belong to this region will be retrieved if the instance_id is `""`.
* `name` - (Optional) The name of MongoDB, all the MongoDBs belong to this region will be retrieved if the name is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `vip` - (Optional) The vip of instances.
* `vnet_id` - (Optional) The ID of subnet. the instance will use the subnet in the current region.
* `vpc_id` - (Optional) The ID of VPC. the instance will use the VPC in the current region.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - It is a nested type which documented below.
  * `area` - availability zone.
  * `config` - instance specification.
  * `create_date` - creation time of the MongoDB.
  * `expiration_date` - expiration date of the MongoDB.
  * `iam_project_id` - ID of the project.
  * `iam_project_name` - Name of the project.
  * `instance_class` - instance specification.
  * `instance_id` - ID of the MongoDB.
  * `instance_type` - Type of the MongoDB.
  * `ip` - A list of MongoDB node IPs.
  * `mode` - MongoDB cluster mode.
  * `mongos_num` - number of mongos.
  * `name` - Name of the MongoDB.
  * `network_type` - Type of network.
  * `node_num` - number of nodes.
  * `pay_type` - type of pay.
  * `port` - Port of the MongoDB.
  * `product_id` - ID of the product.
  * `product_what` - whether the instance is trial or not.
  * `region` - Region.
  * `security_group_id` - ID of the security group.
  * `shard_num` - number of shards.
  * `status` - Status of the MongoDB.
  * `storage` - storage size of the instance disk.
  * `time_cycle` - time cycle of backup.
  * `timezone` - timezone of backup.
  * `timing_switch` - timing switch for backup.
  * `user_id` - ID of the user.
  * `version` - Version of the MongoDB.
  * `vnet_id` - The id of subnet linked to the instance.
  * `vpc_id` - The id of VPC linked to the instance.
* `total_count` - Total number of MongoDBs that satisfy the condition.


