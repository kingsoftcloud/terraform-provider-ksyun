---
layout: "ksyun"
page_title: "Ksyun: ksyun_redis_instance"
sidebar_current: "docs-ksyun-resource-redis-instance"
description: |-
  Provides an Redis instance resource.
---

# ksyun_redis_instance

Provides an redis instance resource.

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
  subnet_name      = "${var.subnet_name}"
  cidr_block = "10.1.0.0/21"
  subnet_type = "Normal"
  dhcp_ip_from = "10.1.0.2"
  dhcp_ip_to = "10.1.0.253"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.1.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  available_zone = "${var.available_zone}"
}

resource "ksyun_redis_sec_group" "default" {
  available_zone = "${var.available_zone}"
  name = "testTerraform777"
  description = "testTerraform777"
}

resource "ksyun_redis_instance" "default" {
  available_zone        = "${var.available_zone}"
  name                  = "MyRedisInstance1101"
  mode                  = 2
  capacity              = 1
  slave_num              = 2  
  net_type              = 2
  vnet_id               = "${ksyun_subnet.default.id}"
  vpc_id                = "${ksyun_vpc.default.id}"
  security_group_id     = "${ksyun_redis_sec_group.default.id}"
  bill_type             = 5
  duration              = ""
  duration_unit         = ""
  pass_word             = "Shiwo1101"
  iam_project_id        = "0"
  protocol              = "${var.protocol}"
  reset_all_parameters  = false
  timing_switch         = "On"
  timezone              = "07:00-08:00"
  available_zone        = "cn-beijing-6a"
  prepare_az_name       = "cn-beijing-6b"
  rr_az_name            = "cn-beijing-6a"
  parameters = {
    "appendonly"                  = "no",
    "appendfsync"                 = "everysec",
    "maxmemory-policy"            = "volatile-lru",
    "hash-max-ziplist-entries"    = "513",
    "zset-max-ziplist-entries"    = "129",
    "list-max-ziplist-size"       = "-2",
    "hash-max-ziplist-value"      = "64",
    "notify-keyspace-events"      = "",
    "zset-max-ziplist-value"      = "64",
    "maxmemory-samples"           = "5",
    "set-max-intset-entries"      = "512",
    "timeout"                     = "600",
  }
}
```

## Argument Reference

The following arguments are supported:

* `available_zone` - (Optional) The Zone to launch the DB instance.
* `name ` - (Optional) The name of DB instance.
* `mode ` - (Optional) The KVStore instance system architecture required by the user. Valid values:  1(cluster),2(single),3(SelfDefineCluster).
* `security_group_id` - (Require) The id of security group;
* `capacity ` - (Require) The instance capacity required by the user. Valid values :{1, 2, 4, 8, 16,20,24,28, 32, 64}.
* `slave_num ` - (Optional) The readonly node num required by the user. Valid values ：{0-7}
* `net_type ` - (Require) The network type. Valid values ：2(vpc).
* `vpc_id` - (Require)   Used to retrieve instances belong to specified VPC .
* `vnet_id` - (Require) The ID of subnet. the instance will use the subnet in the current region.
* `bill_type` - (Optional)Valid values are 1 (Monthly), 5(Daily), 87(HourlyInstantSettlement).
* `duration` - (Optional)Only meaningful if bill_type is 1。 Valid values：{1~36}.
* `duration_unit` - (Optional)Only meaningful if bill_type is 1。 Valid values：month.
* `pass_word` - (Optional)The password of the  instance.The password is a string of 8 to 30 characters and must contain uppercase letters, lowercase letters, and numbers.
* `iam_project_id` - (Optional) The project instance belongs to.
* `protocol` - Engine version. Supported values: 2.8, 4.0 and 5.0.
* `parameters` - Set of parameters needs to be set after instance was launched. Available parameters can refer to the  docs https://docs.ksyun.com/documents/1018 .
* `timing_switch` - (Optional) Switch auto backup. Valid values: On, Off.
* `timezone` - (Optional) Auto backup time zone. Example: "03:00-04:00".
* `shard_size` - (Optional) Shard memory size. If mode is 3 this param is Required.
* `shard_num` - (Optional) Shard num. If mode is 3 this param is Required.
* `prepare_az_name` - (Optional) Assign prepare redis instance az. Mode is 2 this param take effect.
* `rr_az_name` - (Optional) Assign read only redis instance az. Mode is 2 this param take effect.


