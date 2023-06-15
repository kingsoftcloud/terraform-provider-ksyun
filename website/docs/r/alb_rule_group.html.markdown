---
subcategory: "ALB"
layout: "ksyun"
page_title: "ksyun: ksyun_alb_rule_group"
sidebar_current: "docs-ksyun-resource-alb_rule_group"
description: |-
  Provides a ALB rule group resource.
---

# ksyun_alb_rule_group

Provides a ALB rule group resource.

#

## Example Usage

```hcl

```

## Argument Reference

The following arguments are supported:

* `alb_listener_id` - (Required, ForceNew) The ID of the ALB listener.
* `alb_rule_set` - (Required) Rule set.
* `backend_server_group_id` - (Required) The ID of the backend server group.
* `listener_sync` - (Required) Whether to synchronize the health check, session persistence, and load balancing algorithm of the listener. valid values: 'on', 'off'.
* `alb_rule_group_name` - (Optional) The name of the ALB rule group.
* `cookie_name` - (Optional) The name of cookie.
* `cookie_type` - (Optional) The type of cookie, valid values: 'ImplantCookie','RewriteCookie'.
* `health_check_state` - (Optional) Status maintained by health examination.Valid Values:'start', 'stop'.
* `healthy_threshold` - (Optional) Health threshold.Valid Values:1-10. Default is 5.
* `host_name` - (Optional) The service host name of the health check, which is available only for the HTTP or HTTPS health check.
* `interval` - (Optional) Interval of health examination.Valid Values:1-3600. Default is 5.
* `method` - (Optional) Forwarding mode of listener. Valid Values:'RoundRobin', 'LeastConnections'.
* `session_persistence_period` - (Optional) Session hold timeout. Valid Values:1-86400.
* `session_state` - (Optional) The state of session. Valid Values:'start', 'stop'.
* `timeout` - (Optional) Health check timeout.Valid Values:1-3600. Default is 4.
* `unhealthy_threshold` - (Optional) Unhealthy threshold.Valid Values:1-10. Default is 4.
* `url_path` - (Optional) Link to HTTP type listener health check.

The `alb_rule_set` object supports the following:

* `alb_rule_type` - (Required) Rule type. valid values: 'domain', 'url'.
* `alb_rule_value` - (Required) Rule value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `alb_rule_group_id` - The ID of the rule group.


## Import

BWS can be imported using the id, e.g.

```
$ terraform import ksyun_alb_rule_group.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

