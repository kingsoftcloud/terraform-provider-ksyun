---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_listener_associate_acl"
sidebar_current: "docs-ksyun-resource-lb_listener_associate_acl"
description: |-
  Associate a Load Balancer Listener resource with acl.
---

# ksyun_lb_listener_associate_acl

Associate a Load Balancer Listener resource with acl.

#

## Example Usage

```hcl
resource "ksyun_lb_listener_associate_acl" "default" {
  listener_id          = "b330eae5-11a3-4e9e-bf7d-xxxxxxxxxxxx"
  load_balancer_acl_id = "7e94fa82-05c7-496c-ae5e-xxxxxxxxxxxx"
}
```

, e.g.

```hcl
$ terraform import ksyun_lb_listener_associate_acl.default $ { listener_id } : $ { load_balancer_acl_id }
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, ForceNew) The ID of the listener.
* `load_balancer_acl_id` - (Required, ForceNew) The ID of the load balancer acl.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `lb_type` - The type of listener. Valid Value: `Alb` and `Slb`. Default: `Slb`.


