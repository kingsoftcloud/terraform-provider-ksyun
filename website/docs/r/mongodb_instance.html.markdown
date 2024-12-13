---
subcategory: "MongoDB"
layout: "ksyun"
page_title: "ksyun: ksyun_mongodb_instance"
sidebar_current: "docs-ksyun-resource-mongodb_instance"
description: |-
  Provides an replica set MongoDB resource.
---

# ksyun_mongodb_instance

Provides an replica set MongoDB resource.

#

## Example Usage

```hcl
resource "ksyun_mongodb_instance" "default" {
  name              = "InstanceName"
  instance_account  = "root"
  instance_password = "admin"
  instance_class    = "1C2G"
  storage           = 5
  node_num          = 3
  vpc_id            = "VpcId"
  vnet_id           = "VnetId"
  db_version        = "3.6"
  pay_type          = "byDay"
  iam_project_id    = "0"
  availability_zone = "cn-shanghai-3b"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) Availability zone where instance is located.
* `instance_password` - (Required) The administrator password of instance.
* `name` - (Required) The name of instance, which contains 6-64 characters and only support Chinese, English, numbers, '-', '_'.
* `vnet_id` - (Required, ForceNew) The id of subnet linked to the instance.
* `vpc_id` - (Required, ForceNew) The id of VPC linked to the instance.
* `cidrs` - (Optional) network cidr.
* `db_version` - (Optional, ForceNew) The version of instance engine, and support `3.2`, `3.6`, `4.0`, default is `3.2`.
* `duration` - (Optional, ForceNew) The duration of instance use, if `pay_type` is `byMonth`, the duration is required.
* `iam_project_id` - (Optional) The project id of instance belong, if not defined `iam_project_id`, the instance will use `0`.
* `instance_account` - (Optional, ForceNew) The administrator name of instance, if not defined `instance_account`, the instance will use `root`.
* `instance_class` - (Optional) The class of instance cpu and memory.
* `network_type` - (Optional, ForceNew) the type of network.
* `node_num` - (Optional) The num of instance node.
* `pay_type` - (Optional, ForceNew) Instance charge type, if not defined `pay_type`, the instance will use `byMonth`.
* `storage` - (Optional) The size of instance disk, measured in GB (GigaByte).
* `tags` - (Optional) the tags of the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `config` - instance specification.
* `create_date` - creation time of the MongoDB.
* `expiration_date` - expiration date of the MongoDB.
* `iam_project_name` - Name of the project.
* `instance_id` - The id of instance.
* `instance_type` - The type of instance.
* `ip` - IP address.
* `mode` - MongoDB cluster mode.
* `mongos_num` - number of mongos.
* `port` - port number.
* `product_id` - ID of the product.
* `product_what` - whether the instance is trial or not.
* `region` - Region.
* `security_group_id` - The ID of security group.
* `shard_num` - number of shards.
* `status` - the status of instance.
* `time_cycle` - time cycle of backup.
* `timezone` - timezone of backup.
* `timing_switch` - timing switch for backup.
* `user_id` - User ID.
* `version` - Version.


## Import

MongoDB can be imported using the id, e.g.

```
$ terraform import ksyun_mongodb_instance.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

