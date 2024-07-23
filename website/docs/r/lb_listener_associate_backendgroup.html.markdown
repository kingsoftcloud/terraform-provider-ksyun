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

#

## Example Usage

```hcl
resource "ksyun_vpc" "test" {
  vpc_name   = "tf-alb-vpc-1"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "test" {
  subnet_name       = "tf-alb-subnet"
  cidr_block        = "10.0.1.0/24"
  subnet_type       = "Normal"
  availability_zone = var.az
  vpc_id            = ksyun_vpc.test.id
}

resource "ksyun_lb" "default" {
  vpc_id             = ksyun_vpc.test.id
  load_balancer_name = "tf-xun1"
  type               = "public"
}

data "ksyun_certificates" "default" {
}

resource "ksyun_lb_listener" "default" {
  listener_name     = "tf-xun"
  listener_port     = "8000"
  listener_protocol = "TCP"
  listener_state    = "start"
  load_balancer_id  = ksyun_lb.default.id
  method            = "LeastConnections"
  bind_type         = "BackendServerGroup"
  certificate_id    = data.ksyun_certificates.default.certificates.0.certificate_id
}

resource "ksyun_lb_backend_server_group" "default" {
  backend_server_group_name = "xuan-tf"
  vpc_id                    = ksyun_vpc.test.id
  backend_server_group_type = "Server"
  protocol                  = "TCP"
  health_check {
    host_name           = "www.ksyun.com"
    healthy_threshold   = 10
    interval            = 100
    timeout             = 300
    unhealthy_threshold = 10
  }
}

# associate backend server group with listener
resource "ksyun_lb_listener_associate_backendgroup" "mount" {
  listener_id             = ksyun_lb_listener.default.id
  backend_server_group_id = ksyun_lb_backend_server_group.default.id
}
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

