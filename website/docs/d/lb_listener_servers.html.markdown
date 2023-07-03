---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_listener_servers"
sidebar_current: "docs-ksyun-datasource-lb_listener_servers"
description: |-
  Provides a Load Balancer Listener server resource.
---

# ksyun_lb_listener_servers

Provides a Load Balancer Listener server resource.

#

## Example Usage

```hcl
resource "ksyun_lb_listener_server" "default" {
  listener_id      = "3a520244-ddc1-41c8-9d2b-66b4cf3a2386"
  real_server_ip   = "10.0.77.20"
  real_server_port = 8000
  real_server_type = "host"
  instance_id      = "3a520244-ddc1-41c8-9d2b-66b4cf3a2386"
  weight           = 10
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of LB Listener Server IDs, all the LB Listener Servers belong to this region will be retrieved if the ID is `""`.
* `listener_id` - (Optional) A list of LB Listener IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `real_server_ip` - (Optional) A list of real servers.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `servers` - An information list of real servers. Each element contains the following attributes:
  * `instance_id` - Instance ID of the real server, if real server type is host.
  * `master_slave_type` - whether real server is master or salve. when listener method is MasterSlave, this field is supported.
  * `real_server_ip` - IP of the real server.
  * `real_server_port` - Port of the real server.
  * `real_server_state` - State of the real server.
  * `real_server_type` - Type of the real server.
  * `register_id` - register ID of the real server.
  * `weight` - Weight of the real server.
* `total_count` - Total number of LB Listener Servers that satisfy the condition.


