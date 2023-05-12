---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_listeners"
sidebar_current: "docs-ksyun-datasource-listeners"
description: |-
  This data source provides a list of Load Balancer Listener resources according to their Load Balancer Listener ID.
---

# ksyun_listeners

This data source provides a list of Load Balancer Listener resources according to their Load Balancer Listener ID.

#

## Example Usage

```hcl
data "ksyun_listeners" "default" {
  output_file      = "output_result"
  ids              = [""]
  load_balancer_id = ["d3fd0421-a35a-4ddb-a939-5c51e8af8e8c", "4534d617-9de0-4a4a-9ed5-3561196cacb6"]
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Optional) A list of certificate IDs.
* `ids` - (Optional) A list of LB Listener IDs, all the LB Listeners belong to this region will be retrieved if the ID is `""`.
* `load_balancer_id` - (Optional) A list of load balancer IDs.
* `name_regex` - (Optional) A regex string to filter resulting lb listeners by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `listeners` - It is a nested type which documented below.
  * `certificate_id` - ID of the certificate.
  * `create_time` - The time of creation.
  * `enable_http2` - whether support HTTP2.
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
  * `http_protocol` - HTTP protocol.
  * `listener_id` - ID of the LB listener.
  * `listener_name` - Name of the LB listener.
  * `listener_port` - Port of the LB listener.
  * `listener_protocol` - Protocol of the LB listener.
  * `listener_state` - State of the LB listener.
  * `load_balancer_acl_id` - ID of the LB ACL ID.
  * `load_balancer_id` - ID of the LB.
  * `method` - The load balancer method in which the listener is.
  * `real_server` - An information list of real servers. Each element contains the following attributes:
    * `instance_id` - Instance ID of the real server, if real server type is host.
    * `listener_id` - ID of the LB listener.
    * `real_server_ip` - IP of the real server.
    * `real_server_port` - Port of the real server.
    * `real_server_state` - State of the real server.
    * `real_server_type` - Type of the real server.
    * `register_id` - register ID of the real server.
    * `weight` - Weight of the real server.
  * `session` - session configuration.
    * `cookie_name` - The name of cookie.
    * `cookie_type` - The type of the cookie.
    * `session_persistence_period` - Session hold timeout.
    * `session_state` - The state of session.
  * `tls_cipher_policy` - Https listener TLS cipher policy.
* `total_count` - Total number of LB listeners that satisfy the condition.


