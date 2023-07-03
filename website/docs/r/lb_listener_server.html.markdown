---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_listener_server"
sidebar_current: "docs-ksyun-resource-lb_listener_server"
description: |-
  Provides a Load Balancer Listener server resource.
---

# ksyun_lb_listener_server

Provides a Load Balancer Listener server resource.

#

## Example Usage

```hcl
resource "ksyun_lb_listener_server" "default" {
  listener_id      = "3a520244-ddc1-41c8-9d2b-xxxxxxxxxxxx"
  real_server_ip   = "10.0.77.20"
  real_server_port = 8000
  real_server_type = "host"
  instance_id      = "3a520244-ddc1-41c8-9d2b-xxxxxxxxxxxx"
  weight           = 10
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, ForceNew) The id of the listener.
* `real_server_ip` - (Required, ForceNew) The IP of real server.
* `real_server_port` - (Required) The port of real server.Valid Values:1-65535.
* `instance_id` - (Optional, ForceNew) The ID of instance.
* `master_slave_type` - (Optional) whether real server is master of salve. when listener method is MasterSlave, this field is supported.
* `real_server_type` - (Optional, ForceNew) The type of real server.Valid Values:'host', 'DirectConnectGateway', 'VpnTunnel'.
* `weight` - (Optional) The weight of backend service.Valid Values:1-255.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `listener_method` - Forwarding mode of listener.
* `real_server_state` - State of the real server.
* `register_id` - The registration ID of real server.


## Import

LB Listener can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb_listener.example vserver-abcdefg
```

