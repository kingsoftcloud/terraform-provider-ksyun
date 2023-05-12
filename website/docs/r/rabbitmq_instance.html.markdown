---
subcategory: "RabbitMQ"
layout: "ksyun"
page_title: "ksyun: ksyun_rabbitmq_instance"
sidebar_current: "docs-ksyun-resource-rabbitmq_instance"
description: |-
  Provides an replica set Rabbitmq resource.
---

# ksyun_rabbitmq_instance

Provides an replica set Rabbitmq resource.

#

## Example Usage

```hcl
resource "ksyun_rabbitmq_instance" "default" {
  availability_zone = "cn-beijing-6a"
  instance_name     = "my_rabbitmq_instance"
  instance_password = "Shiwo1101"
  instance_type     = "2C4G"
  vpc_id            = "VpcId"
  subnet_id         = "VnetId"
  mode              = 1
  engine_version    = "3.7"
  ssd_disk          = "5"
  node_num          = 3
  bill_type         = 87
  project_id        = 103800
}
```

## Argument Reference

The following arguments are supported:

* `bill_type` - (Required, ForceNew) Instance charge type,Valid values are 1 (Monthly), 87(UsageInstantSettlement).
* `engine_version` - (Required, ForceNew) The version of instance engine.
* `instance_name` - (Required) The name of instance, which contains 6-64 characters and only support Chinese, English, numbers, '-', '_'.
* `instance_password` - (Required) The administrator password of instance.
* `instance_type` - (Required, ForceNew) The class of instance cpu and memory.
* `mode` - (Required, ForceNew) The mode of instance.
* `ssd_disk` - (Required, ForceNew) The size of instance disk, measured in GB (GigaByte).
* `subnet_id` - (Required, ForceNew) The id of subnet linked to the instance.
* `vpc_id` - (Required, ForceNew) The id of VPC linked to the instance.
* `availability_zone` - (Optional, ForceNew) Availability zone where instance is located.
* `cidrs` - (Optional) network cidrs.
* `duration` - (Optional, ForceNew) The duration of instance use, if `bill_type` is `1`, the duration is required.
* `enable_eip` - (Optional) If the value is true, the instance will support public ip. default is false.
* `enable_plugins` - (Optional) Enable plugins.
* `force_restart` - (Optional) Set it to true to make some parameter efficient when modifying them. Default to false.
* `node_num` - (Optional, ForceNew) the number of instance node, if not defined 'node_num', the instance will use '3'.
* `project_id` - (Optional) The project id of instance belong, if not defined `project_id`, the instance will use `0`.
* `security_group_id` - (Optional) security group id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_date` - The creation date.
* `eip_egress` - The egress of the EIP.
* `eip` - EIP address.
* `engine` - engine.
* `expiration_date` - The expiration date.
* `instance_id` - The ID of the instance.
* `mode_name` - mode name.
* `network_type` - Network type.
* `port` - The port of the instance.
* `product_id` - The ID of the project.
* `product_what` - product what.
* `project_name` - The name of the project.
* `protocol` - protocol.
* `region` - region.
* `status_name` - status name.
* `status` - The status of the instance.
* `user_id` - user id.
* `vip` - vip.
* `web_eip` - Web EIP address.
* `web_vip` - web vip.


