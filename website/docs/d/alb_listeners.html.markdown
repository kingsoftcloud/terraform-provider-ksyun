---
subcategory: "ALB"
layout: "ksyun"
page_title: "ksyun: ksyun_alb_listeners"
sidebar_current: "docs-ksyun-datasource-alb_listeners"
description: |-
  This data source provides a list of ALB listener resources according to their ID.
---

# ksyun_alb_listeners

This data source provides a list of ALB listener resources according to their ID.

#

## Example Usage

```hcl
data "ksyun_alb_listeners" "default" {
  output_file = "output_result"
  ids         = []
  alb_id      = []
  acl_id      = []
  protocol    = []
}
```

## Argument Reference

The following arguments are supported:

* `acl_id` - (Optional) One or more ACL ID.
* `alb_id` - (Optional) One or more ALB IDs.
* `ids` - (Optional) A list of ALB Listener IDs, all the ALB Listeners belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `protocol` - (Optional) One or more Listener protocol.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `listeners` - An information list of ALB Listeners. Each element contains the following attributes:
  * `alb_id` - The ID of the ALB Listener.
  * `alb_listener_id` - The ID of the listener.
  * `alb_listener_name` - The name of the listener.
  * `alb_listener_state` - The state of listener.
  * `certificate_id` - The ID of certificate.
  * `create_time` - The creation time.
  * `enable_http2` - whether enable to HTTP2.
  * `health_check` - Health check configuration. It is a nested type which documented below.
    * `health_check_id` - ID of the healthcheck.
    * `health_check_state` - Status maintained by health examination.
    * `healthy_threshold` - Health threshold.
    * `host_name` - The service host name of the health check, which is available only for the HTTP or HTTPS health check.
    * `interval` - Interval of health examination.
    * `listener_id` - ID of the LB listener.
    * `timeout` - Health check timeout.
    * `unhealthy_threshold` - Unhealthy threshold.
    * `url_path` - Link to HTTP type listener health check.
  * `http_protocol` - Backend Protocol.
  * `id` - ID of the ALB Listener.
  * `method` - Forwarding mode of listener.
  * `port` - The protocol port of listener.
  * `protocol` - The protocol of listener.
  * `redirect_alb_listener_id` - The ID of the redirect ALB listener.
  * `session` - session.
    * `cookie_name` - The name of cookie.
    * `cookie_type` - The type of cookie.
    * `session_persistence_period` - Session hold timeout.
    * `session_state` - The state of session.
  * `tls_cipher_policy` - TLS cipher policy.
* `total_count` - Total number of ALB listeners that satisfy the condition.


