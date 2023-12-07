---
subcategory: "ALB"
layout: "ksyun"
page_title: "ksyun: ksyun_alb_register_backend_server"
sidebar_current: "docs-ksyun-resource-alb_register_backend_server"
description: |-
  Provides alb register alb backend server group resource.
---

# ksyun_alb_register_backend_server

Provides alb register alb backend server group resource.

#

## Example Usage

```hcl
resource "ksyun_vpc" "test" {
  vpc_name   = "tf-alb-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_alb_backend_server_group" "foo" {
  name                = "tf-alb-bsg"
  vpc_id              = ksyun_vpc.test.id
  upstream_keepalive  = "adaptation"
  backend_server_type = "Host"
}

resource "ksyun_alb_register_backend_server" "foo" {
  backend_server_group_id = ksyun_alb_backend_server_group.foo.id
  backend_server_ip       = ksyun_instance.test.private_ip_address
  port                    = 8080
  weight                  = 40
}
```

## Argument Reference

The following arguments are supported:

* `backend_server_group_id` - (Required, ForceNew) The ID of alb backend server group.
* `backend_server_ip` - (Required, ForceNew) The IP of alb backend server.
* `port` - (Required, ForceNew) The port of alb backend server. Valid Values:1-65535.
* `direct_connect_gateway_id` - (Optional, ForceNew) The ID of direct connect gateway.
* `network_interface_id` - (Optional, ForceNew) The ID of network interface.
* `weight` - (Optional) The weight of backend service. Valid Values:0-255.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `backend_server_id` - The registration ID of binding server group.
* `create_time` - The time when the backend service was created.
* `instance_id` - The ID of instance.


## Import

resource can be imported using the id, e.g.

```
$ terraform import ksyun_alb_register_backend_server.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

