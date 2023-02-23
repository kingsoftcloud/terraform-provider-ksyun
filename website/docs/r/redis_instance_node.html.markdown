---
subcategory: "Redis"
layout: "ksyun"
page_title: "ksyun: ksyun_redis_instance_node"
sidebar_current: "docs-ksyun-resource-redis_instance_node"
description: |-
  Provides an redis instance node resource.
---

# ksyun_redis_instance_node

Provides an redis instance node resource.

#

## Example Usage

```hcl
resource "ksyun_redis_instance_node" "default" {
  cache_id       = "${ksyun_redis_instance.default.id}"
  available_zone = "${var.available_zone}"
}

resource "ksyun_redis_instance_node" "node" {
  // creating multiple read-only nodes,
  // not concurrently, requires dependencies to synchronize the execution of creating multiple read-only nodes.
  // if only one read-only node is created, it is not required to fill in.
  pre_node_id    = "${ksyun_redis_instance_node.default.id}"
  cache_id       = "${ksyun_redis_instance.default.id}"
  available_zone = "${var.available_zone}"
}
```

## Argument Reference

The following arguments are supported:

* `cache_id` - (Required, ForceNew) The ID of the instance.
* `available_zone` - (Optional, ForceNew) The Zone to launch the DB instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - creation time.
* `instance_id` - The ID of the instance.
* `ip` - IP address.
* `name` - Name.
* `port` - Port number.
* `proxy` - proxy.
* `status` - status.


## Import

redis node can be imported using the `id`, e.g.

```
$ terraform import ksyun_redis_instance_node.default xxxxxxxxx
```

