---
subcategory: "ALB"
layout: "ksyun"
page_title: "ksyun: ksyun_alb_rule_groups"
sidebar_current: "docs-ksyun-datasource-alb_rule_groups"
description: |-
  This data source provides a list of ALB rule group resources according to their ID.
---

# ksyun_alb_rule_groups

This data source provides a list of ALB rule group resources according to their ID.

#

## Example Usage

```hcl
data "ksyun_alb_rule_groups" "default" {
  output_file     = "output_result"
  ids             = []
  alb_listener_id = []
}
```

## Argument Reference

The following arguments are supported:

* `alb_listener_id` - (Optional) one or more alb listener id.
* `ids` - (Optional) A list of ALB Rule Group IDs, all the ALB Rule Group belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `alb_rule_groups` - An information list of ALB Rule Groups. Each element contains the following attributes:
  * `alb_listener_id` - The ID of the ALB listener.
  * `alb_rule_group_id` - The ID of the rule group.
  * `alb_rule_group_name` - The name of the ALB rule group.
  * `alb_rule_set` - Rule set.
    * `alb_rule_type` - Rule type.
    * `alb_rule_value` - Rule value.
  * `backend_server_group_id` - The ID of the backend server group.
  * `cookie_name` - The name of cookie.
  * `cookie_type` - The type of cookie.
  * `health_check_state` - Status maintained by health examination.
  * `health_port` - The port of connecting for health check.
  * `health_protocol` - The protocol of connecting for health check.
  * `healthy_threshold` - Health threshold.
  * `host_name` - The service host name of the health check.
  * `http_method` - The http requests' method.
  * `id` - ID of the ALB Rule Group.
  * `interval` - Interval of health examination.
  * `listener_sync` - Whether to synchronize the health check, session persistence, and load balancing algorithm of the listener.
  * `method` - Forwarding mode of listener.
  * `session_persistence_period` - Session hold timeout.
  * `session_state` - The state of session.
  * `timeout` - Health check timeout.
  * `unhealthy_threshold` - Unhealthy threshold.
  * `url_path` - Link to HTTP type listener health check.
* `total_count` - Total number of ALB Rule Groups that satisfy the condition.


