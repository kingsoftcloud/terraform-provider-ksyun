---
subcategory: "ALB"
layout: "ksyun"
page_title: "ksyun: ksyun_alb_backend_server_groups"
sidebar_current: "docs-ksyun-datasource-alb_backend_server_groups"
description: |-
  Provides a list of lb AlbBackend server groups in the current region.
---

# ksyun_alb_backend_server_groups

Provides a list of lb AlbBackend server groups in the current region.

#

## Example Usage

```hcl
# Get availability zones
data "ksyun_alb_backend_server_groups" "default" {
  output_file = "out_file"
  ids         = []
}
```

## Argument Reference

The following arguments are supported:

* `alb_backend_server_group_type` - (Optional) A list of AlbBackendServerGroup types.
* `ids` - (Optional) A list of AlbBackendServerGroup IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `vpc_id` - (Optional) A list of VPC IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `alb_backend_server_groups` - An information list of AlbBackendServerGroups. Each element contains the following attributes:
  * `backend_server_group_id` - The id of AlbBackend server group.
  * `backend_server_group_type` - The type of blb backend server group.
  * `backend_server_number` - The number of alb backend server number.
  * `create_time` - The time when the alb backend server group was created.
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
  * `name` - The name of AlbBackend server group.
  * `upstream_keepalive` - The upstream keepalive type.
  * `vpc_id` - Virtual private network ID.
* `total_count` - Total number of AlbBackendServerGroups that satisfy the condition.


