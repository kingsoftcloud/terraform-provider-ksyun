---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_rules"
sidebar_current: "docs-ksyun-datasource-lb_rules"
description: |-
  Provides a list of ksyun lb rules resources in the current region.
---

# ksyun_lb_rules

Provides a list of ksyun lb rules resources in the current region.

#

## Example Usage

```hcl
data "ksyun_lb_rules" "default" {
  output_file    = "output_result"
  ids            = []
  host_header_id = []
}
```

## Argument Reference

The following arguments are supported:

* `host_header_id` - (Optional) The id of host header.
* `ids` - (Optional) A list of rule IDs.
* `output_file` - (Optional) File name where to save data source results (after running terraform plan).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `lb_rules` - An information list of LB Rules. Each element contains the following attributes:
  * `backend_server_group_id` - ID of the backend server group.
  * `create_time` - The time of creation for LB Rule.
  * `health_check` - A list of health checks.
    * `health_check_id` - ID of the health check.
    * `health_check_state` - Status maintained by health examination.
    * `healthy_threshold` - Health threshold.
    * `host_name` - Domain name of HTTP type health check.
    * `interval` - Interval of health examination.
    * `listener_id` - ID of the listener.
    * `timeout` - Health check timeout.
    * `unhealthy_threshold` - Unhealthy threshold.
    * `url_path` - Link to HTTP type listener health check.
  * `host_header_id` - ID of the host header.
  * `listener_sync` - Whether to synchronizethe the health check, the session hold and the forward algorithms of the listener.Valid Values:'on', 'off'.
  * `method` - Forwarding mode of listener.
  * `path` - The path of rule.
  * `rule_id` - ID of the rule.
  * `session` - session configuration.
    * `cookie_name` - The name of cookie.
    * `cookie_type` - The type of the cookie.
    * `session_persistence_period` - Session hold timeout.
    * `session_state` - The state of session.
* `total_count` - Total number of LB Rules that satisfy the condition.


