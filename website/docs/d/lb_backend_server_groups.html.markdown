---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_backend_server_groups"
sidebar_current: "docs-ksyun-datasource-lb_backend_server_groups"
description: |-
  Provides a list of lb backend server groups in the current region.
---

# ksyun_lb_backend_server_groups

Provides a list of lb backend server groups in the current region.

#

## Example Usage

```hcl
# Get availability zones
data "ksyun_lb_backend_server_groups" "default" {
  output_file = "out_file"
  ids         = []
}
```

## Argument Reference

The following arguments are supported:

* `backend_server_group_type` - (Optional) A list of BackendServerGroup types.
* `ids` - (Optional) A list of BackendServerGroup IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `vpc_id` - (Optional) A list of VPC IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backend_server_groups` - An information list of BackendServerGroups. Each element contains the following attributes:
  * `backend_server_group_id` - The id of backend server group.
  * `backend_server_group_name` - The name of backend server group.
  * `backend_server_group_type` - The type of backend server group.Valid values are Server and Mirror.
  * `backend_server_number` - The number of backend server number.
  * `create_time` - The time when the backend server group was created.
  * `health_check` - Health check information, only the mirror server has this parameter.
    * `health_check_id` - ID of the health check.
    * `health_check_state` - state of the health check.
    * `healthy_threshold` - health threshold.
    * `host_name` - hostname of the health check.
    * `interval` - interval of the health check.
    * `listener_id` - ID of the listener.
    * `timeout` - timeout of the health check.
    * `unhealthy_threshold` - Unhealthy threshold of health check.
    * `url_path` - path of the health check.
  * `protocol` - The protocol of the backend server group. Valid values: 'TCP', 'UDP', 'HTTP'. Default `HTTP`.
  * `vpc_id` - Virtual private network ID.
* `total_count` - Total number of BackendServerGroups that satisfy the condition.


