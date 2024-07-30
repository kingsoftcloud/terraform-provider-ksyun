---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_listener_associate_backendgroup"
sidebar_current: "docs-ksyun-resource-lb_listener_associate_backendgroup"
description: |-
  Provides slb listener mount backend server group resource.
---

# ksyun_lb_listener_associate_backendgroup

Provides slb listener mount backend server group resource.

~> **NOTE:** This resource is **deprecated**. Use `backend_server_group_mounted` of `ksyun_lb_listener` instead. See [ksyun_lb_listener](https://registry.terraform.io/providers/kingsoftcloud/ksyun/latest/docs/resources/lb_listener) for more details.

#

## Example Usage

```hcl

```

## Argument Reference

The following arguments are supported:

* `backend_server_group_id` - (Required, ForceNew) The ID of alb backend server group.
* `listener_id` - (Required, ForceNew) The ID of slb listener.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

resource can be imported using the id, e.g.

```
$ terraform import ksyun_lb_listener_associate_backendgroup.default $listener_id:$backend_server_group_id
```

