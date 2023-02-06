---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_register_backend_servers"
sidebar_current: "docs-ksyun-datasource-lb_register_backend_servers"
description: |-
  Provides a list of register backend servers in the current region.
---

# ksyun_lb_register_backend_servers

Provides a list of register backend servers in the current region.

#

## Example Usage

```hcl
data "ksyun_lb_register_backend_servers" "foo" {
  output_file             = "output_result"
  ids                     = []
  backend_server_group_id = []
}
```

## Argument Reference

The following arguments are supported:

* `backend_server_group_id` - (Optional) A list of Register backend server IDs.
* `ids` - (Optional) A list of Register backend server IDs, all the Register backend servers belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `register_backend_servers` - An information list of Register backend groups. Each element contains the following attributes:
  * `backend_server_group_id` - The ID of the server.
  * `backend_server_ip` - The IP of the server.
  * `create_time` - creation time.
  * `instance_id` - The ID of the instance.
  * `network_interface_id` - The ID of the network interface.
  * `real_server_ip` - The IP of real server.
  * `real_server_port` - The port number of real server.
  * `real_server_state` - The state of the real server.
  * `real_server_type` - The type of real server.
  * `register_id` - The registration ID of the binding server group.
  * `weight` - The weight of backend service.
* `total_count` - Total number of Register backend groups that satisfy the condition.


