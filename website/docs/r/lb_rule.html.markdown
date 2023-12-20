---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_rule"
sidebar_current: "docs-ksyun-resource-lb_rule"
description: |-
  Provides a lb rule resource.
---

# ksyun_lb_rule

Provides a lb rule resource.

#

## Example Usage

```hcl
resource "ksyun_lb_rule" "default" {
  path                    = "/tfxun/update",
  host_header_id          = "",
  backend_server_group_id = ""
  listener_sync           = "on"
  method                  = "RoundRobin"
  session {
    session_state              = "start"
    session_persistence_period = 1000
    cookie_type                = "ImplantCookie"
    cookie_name                = "cookiexunqq"
  }
  health_check {
    health_check_state  = "start"
    healthy_threshold   = 2
    interval            = 200
    timeout             = 2000
    unhealthy_threshold = 2
    url_path            = "/monitor"
    host_name           = "www.ksyun.com"
  }
}
```

## Argument Reference

The following arguments are supported:

* `backend_server_group_id` - (Required) The id of backend server group.
* `host_header_id` - (Required, ForceNew) The id of host header id.
* `path` - (Required) The path of rule.
* `health_check` - (Optional) health check configuration.
* `listener_sync` - (Optional) Whether to synchronizethe the health check, the session hold and the forward algorithms of the listener.Valid Values:'on', 'off'. Default is 'on'.
* `method` - (Optional) Forwarding mode of listener.Valid Values:'RoundRobin', 'LeastConnections'. Default is 'RoundRobin'.
* `session` - (Optional) Session.

The `health_check` object supports the following:

* `health_check_connect_port` - (Optional) The port of connecting for health check.
* `health_check_state` - (Optional) Status maintained by health examination.Valid Values:'start', 'stop'.
* `healthy_threshold` - (Optional) Health threshold.Valid Values:1-10. Default is 5.
* `host_name` - (Optional) The service host name of the health check, which is available only for the HTTP or HTTPS health check.
* `interval` - (Optional) Interval of health examination.Valid Values:1-3600. Default is 5.
* `is_default_host_name` - (Optional) Whether the host name is default or not.
* `timeout` - (Optional) Health check timeout.Valid Values:1-3600. Default is 4.
* `unhealthy_threshold` - (Optional) Unhealthy threshold.Valid Values:1-10. Default is 4.
* `url_path` - (Optional) Link to HTTP type listener health check.

The `session` object supports the following:

* `cookie_name` - (Optional) The name of cookie.The CookieType is valid and required when it is 'RewriteCookie'; otherwise, this value is ignored.
* `cookie_type` - (Optional) The type of the cookie.Valid Values:'ImplantCookie', 'RewriteCookie'. Default is 'ImplantCookie'.
* `session_persistence_period` - (Optional) Session hold timeout.Valid Values:1-86400. Default is '7200'.
* `session_state` - (Optional) The state of session.Valid Values:'start', 'stop'. Default is 'start'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - the creation time.
* `rule_id` - The ID of the rule.


## Import

LB Rule can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb_rule.example vserver-abcdefg
```

