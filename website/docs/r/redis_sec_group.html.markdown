---
subcategory: "Redis"
layout: "ksyun"
page_title: "ksyun: ksyun_redis_sec_group"
sidebar_current: "docs-ksyun-resource-redis_sec_group"
description: |-
  Provides an redis security group function.
---

# ksyun_redis_sec_group

Provides an redis security group function.

#

## Example Usage

```hcl
variable "available_zone" {
  default = "cn-beijing-6a"
}

resource "ksyun_redis_sec_group" "add" {
  available_zone = "${var.available_zone}"
  name           = "testAddTerraform"
  description    = "testAddTerraform"
}

resource "ksyun_redis_sec_group_rule" "default" {
  available_zone    = "${var.available_zone}"
  security_group_id = "${ksyun_redis_sec_group.add.id}"
  rules             = ["172.16.0.0/32", "192.168.0.0/32"]
}

resource "ksyun_redis_sec_group_allocate" "default" {
  available_zone    = "${var.available_zone}"
  security_group_id = "${ksyun_redis_sec_group.add.id}"
  cache_ids         = ["122334234"]
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required) The description of the security group.
* `name` - (Required) The name of the security group.
* `available_zone` - (Optional, ForceNew) The Zone to launch the security group.
* `cache_ids` - (Optional) The ids of the redis instance.
* `rules` - (Optional) The cidr block of source for the instance, multiple cidr separated by comma.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Redis security group can be imported using the `id`, e.g.

```
$ terraform import ksyun_redis_sec_group.default fdeba8ca-8aa6-4cd0-8ffa-xxxxxxxxxxxx
```

