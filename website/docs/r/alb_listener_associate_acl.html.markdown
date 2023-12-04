---
subcategory: "ALB"
layout: "ksyun"
page_title: "ksyun: ksyun_alb_listener_associate_acl"
sidebar_current: "docs-ksyun-resource-alb_listener_associate_acl"
description: |-
  Associate an Application Load Balancer Listener resource with acl.
---

# ksyun_alb_listener_associate_acl

Associate an Application Load Balancer Listener resource with acl.

#

## Example Usage

```hcl
resource "ksyun_alb_Listener_associate_acl" "default" {
  AlbListener_id       = "b330eae5-11a3-4e9e-bf7d-xxxxxxxxxxxx"
  load_balancer_acl_id = "7e94fa82-05c7-496c-ae5e-xxxxxxxxxxxx"
}
```

, e.g.

```hcl
$ terraform import ksyun_lb_AlbListener_associate_acl.default $ { listener_id } : $ { load_balancer_acl_id }
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, ForceNew) The ID of the AlbListener.
* `load_balancer_acl_id` - (Required, ForceNew) The ID of the load balancer acl.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `lb_type` - The type of listener. Valid Value: `Alb` and `Slb`. Default: `Slb`.


