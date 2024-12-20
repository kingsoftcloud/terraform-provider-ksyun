---
subcategory: "Redis"
layout: "ksyun"
page_title: "ksyun: ksyun_redis_instance"
sidebar_current: "docs-ksyun-resource-redis_instance"
description: |-
  Provides an redis instance resource.
---

# ksyun_redis_instance

Provides an redis instance resource.

#

## Example Usage

```hcl
variable "available_zone" {
  default = "cn-beijing-6a"
}

variable "subnet_name" {
  default = "ksyun_subnet_tf"
}

variable "vpc_name" {
  default = "ksyun_vpc_tf"
}

variable "vpc_cidr" {
  default = "10.1.0.0/21"
}

variable "protocol" {
  default = "4.0"
}

resource "ksyun_vpc" "default" {
  vpc_name   = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}

resource "ksyun_subnet" "default" {
  subnet_name    = "${var.subnet_name}"
  cidr_block     = "10.1.0.0/21"
  subnet_type    = "Normal"
  dhcp_ip_from   = "10.1.0.2"
  dhcp_ip_to     = "10.1.0.253"
  vpc_id         = "${ksyun_vpc.default.id}"
  gateway_ip     = "10.1.0.1"
  dns1           = "198.18.254.41"
  dns2           = "198.18.254.40"
  available_zone = "${var.available_zone}"
}

resource "ksyun_redis_sec_group" "default" {
  available_zone = "${var.available_zone}"
  name           = "testTerraform777"
  description    = "testTerraform777"
}

resource "ksyun_redis_instance" "default" {
  available_zone       = "${var.available_zone}"
  name                 = "MyRedisInstance1101"
  mode                 = 2
  capacity             = 1
  slave_num            = 2
  net_type             = 2
  vnet_id              = "${ksyun_subnet.default.id}"
  vpc_id               = "${ksyun_vpc.default.id}"
  security_group_id    = "${ksyun_redis_sec_group.default.id}"
  bill_type            = 5
  duration             = ""
  duration_unit        = ""
  pass_word            = "Shiwo1101"
  iam_project_id       = "0"
  protocol             = "${var.protocol}"
  reset_all_parameters = false
  timing_switch        = "On"
  timezone             = "07:00-08:00"
  available_zone       = "cn-beijing-6a"
  prepare_az_name      = "cn-beijing-6b"
  rr_az_name           = "cn-beijing-6a"
  parameters = {
    "appendonly"               = "no",
    "appendfsync"              = "everysec",
    "maxmemory-policy"         = "volatile-lru",
    "hash-max-ziplist-entries" = "513",
    "zset-max-ziplist-entries" = "129",
    "list-max-ziplist-size"    = "-2",
    "hash-max-ziplist-value"   = "64",
    "notify-keyspace-events"   = "",
    "zset-max-ziplist-value"   = "64",
    "maxmemory-samples"        = "5",
    "set-max-intset-entries"   = "512",
    "timeout"                  = "600",
  }
}
```

## Argument Reference

The following arguments are supported:

* `capacity` - (Required) mem of redis, if mode is selfDefineCluster(mode=3) then capacity = shard_size * shard_num.
* `name` - (Required) The name of DB instance.
* `vnet_id` - (Required, ForceNew) The ID of subnet. the instance will use the subnet in the current region.
* `vpc_id` - (Required, ForceNew) Used to retrieve instances belong to specified VPC.
* `available_zone` - (Optional, ForceNew) The Zone to launch the DB instance.
* `backup_time_zone` - (Optional) Auto backup time zone. Example: "03:00-04:00".
* `bill_type` - (Optional, ForceNew) Valid values are 1 (Monthly), 5(Daily), 87(HourlyInstantSettlement).
* `delete_directly` - (Optional) Default is `false`, deleted instance will remain in the recycle bin. Setting the value to `true`, instance is permanently deleted without being recycled.
* `duration` - (Optional, ForceNew) Only meaningful if bill_type is 1, Valid values:{1~36}.
* `iam_project_id` - (Optional) The project instance belongs to.
* `mode` - (Optional, ForceNew) The KVStore instance system architecture required by the user. Valid values:  1(cluster),2(single),3(SelfDefineCluster).
* `net_type` - (Optional) The network type. Valid values: 2(vpc).
* `parameters` - (Optional) Set of parameters needs to be set after instance was launched. Available parameters can refer to the  docs https://docs.ksyun.com/documents/1018.
* `pass_word` - (Optional) The password of the  instance.The password is a string of 8 to 30 characters and must contain uppercase letters, lowercase letters, and numbers.
* `prepare_az_name` - (Optional, ForceNew) assign standby instance area.
* `protocol` - (Optional, ForceNew) Engine version. Supported values: 2.8, 4.0 and 5.0.
* `reset_all_parameters` - (Optional) whether reset all parameters.
* `rr_az_name` - (Optional, ForceNew) assign read only instance area.
* `security_group_id` - (Optional) The id of security group.
* `shard_num` - (Optional) shard number.
* `shard_size` - (Optional) each shard mem size GB.
* `slave_num` - (Optional, ForceNew) The readonly node num required by the user. Valid values: {0-7}.
* `tags` - (Optional) the tags of the resource.
* `timezone` - (Optional) Auto backup time zone. Example: "03:00-04:00".
* `timing_switch` - (Optional) auto backup On or Off.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `az` - availability zone.
* `cache_id` - The ID of cache.
* `create_time` - creation time.
* `engine` - engine.
* `iam_project_name` - project name.
* `order_type` - order type.
* `order_use` - order use.
* `port` - port.
* `product_id` - project id.
* `service_begin_time` - service begin time.
* `service_end_time` - service end time.
* `service_status` - service status.
* `size` - size.
* `slave_vip` - slave vip.
* `source` - source.
* `status` - status.
* `sub_order_id` - sub order ID.
* `used_memory` - used memory.
* `vip` - vip.


## Import

redis  can be imported using the `id`, e.g.

```
$ terraform import ksyun_redis_instance.example xxxxxxxxx
```

