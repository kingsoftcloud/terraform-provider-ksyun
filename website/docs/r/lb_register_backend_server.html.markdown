---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_register_backend_server"
sidebar_current: "docs-ksyun-resource-lb_register_backend_server"
description: |-
  Provides a lb register backend server resource.
---

# ksyun_lb_register_backend_server

Provides a lb register backend server resource.

#

## Example Usage

```hcl
resource "ksyun_lb_register_backend_server" "default" {
  backend_server_group_id = "xxxx"
  backend_server_ip       = "192.168.5.xxx"
  backend_server_port     = "8081"
  weight                  = 10
}
```

## Argument Reference

The following arguments are supported:

* `backend_server_group_id` - (Required, ForceNew) The ID of backend server group.
* `backend_server_ip` - (Required, ForceNew) The IP of backend server.
* `backend_server_port` - (Required, ForceNew) The port of backend server.Valid Values:1-65535.
* `weight` - (Optional) The weight of backend service.Valid Values:0-255.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time when the backend service was created.
* `instance_id` - The ID of instance.
* `network_interface_id` - The ID of network interface.
* `real_server_ip` - The IP of real server.
* `real_server_port` - The port of real server.Valid Values:1-65535.
* `real_server_state` - The state of real server.Values:'healthy','unhealthy'.
* `real_server_type` - The type of real server.Valid Values:'Host'.
* `register_id` - The registration ID of binding server group.


## Import

resource can be imported using the id, e.g.

```
$ terraform import ksyun_lb_register_backend_server.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

