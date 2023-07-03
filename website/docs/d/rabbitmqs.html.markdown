---
subcategory: "RabbitMQ"
layout: "ksyun"
page_title: "ksyun: ksyun_rabbitmqs"
sidebar_current: "docs-ksyun-datasource-rabbitmqs"
description: |-
  This data source provides a list of Rabbitmq resources according to their name, Instance ID, Subnet ID, VPC ID and the Project ID they belong to.
---

# ksyun_rabbitmqs

This data source provides a list of Rabbitmq resources according to their name, Instance ID, Subnet ID, VPC ID and the Project ID they belong to.

#

## Example Usage

```hcl
data "ksyun_rabbitmqs" "default" {
  output_file   = "output_result"
  project_id    = ""
  instance_id   = ""
  instance_name = ""
  subnet_id     = ""
  vpc_id        = ""
  vip           = ""
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional) The id of Rabbitmq, all the Rabbitmqs belong to this region will be retrieved if the instance_id is `""`.
* `instance_name` - (Optional) The name of RabbitMQ.
* `name` - (Optional) The name of RabbitMQ.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `project_id` - (Optional) One or more project IDs.
* `subnet_id` - (Optional) The ID of the subnet.
* `total_count` - (Optional) Total number of RabbitMQs that satisfy the condition.
* `vip` - (Optional) The vip of RabbitMQs.
* `vpc_id` - (Optional) The ID of the VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - An information list of RabbitMQ instances. Each element contains the following attributes:
  * `availability_zone` - Availability zone where instance is located.
  * `bill_type` - Instance charge type.
  * `create_date` - Time of creation.
  * `duration` - The duration of instance use.
  * `eip_egress` - egress of the EIP.
  * `eip` - EIP.
  * `engine_version` - The version of instance engine.
  * `engine` - the engine of the instance.
  * `expiration_date` - Time of expiration.
  * `instance_id` - The ID of the RabbitMQ instance.
  * `instance_name` - The name of the RabbitMQ instance.
  * `instance_password` - The administrator password of instance.
  * `instance_type` - The class of instance cpu and memory.
  * `mode_name` - The mode name of the instance.
  * `mode` - The mode of instance.
  * `network_type` - the network type of the instance.
  * `node_num` - the number of instance node.
  * `port` - port number.
  * `product_id` - the product id of the instance.
  * `product_what` - whether the instance is trial or not.
  * `project_id` - The project instance belongs to.
  * `project_name` - The project name of instance belong.
  * `protocol` - the protocol of the instance.
  * `region` - region.
  * `security_group_id` - the id of the security group.
  * `ssd_disk` - The size of instance disk, measured in GB (GigaByte).
  * `status_name` - the status name of the instance.
  * `status` - the status of the instance.
  * `subnet_id` - The id of subnet linked to the instance.
  * `user_id` - the id of user.
  * `vip` - the vip of the instance.
  * `vpc_id` - The id of VPC linked to the instance.
  * `web_eip` - Web EIP.
  * `web_vip` - the web vip of the instance.


