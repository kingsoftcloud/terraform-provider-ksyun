---
subcategory: "ALB"
layout: "ksyun"
page_title: "ksyun: ksyun_alb_listener"
sidebar_current: "docs-ksyun-resource-alb_listener"
description: |-
  Provides a ALB Listener resource.
---

# ksyun_alb_listener

Provides a ALB Listener resource.

#

## Example Usage

```hcl
resource "ksyun_alb_listener" "default" {
}
```

## Argument Reference

The following arguments are supported:

* `alb_id` - (Required, ForceNew) The ID of the ALB.
* `alb_listener_state` - (Required) The state of listener.Valid Values:'start', 'stop'.
* `port` - (Required, ForceNew) The protocol port of listener.
* `protocol` - (Required, ForceNew) The protocol of listener. Valid Values: 'HTTP', 'HTTPS'.
* `alb_listener_name` - (Optional) The name of the listener.
* `certificate_id` - (Optional) The ID of certificate.
* `enable_http2` - (Optional) whether enable to HTTP2.
* `http_protocol` - (Optional) Backend Protocol, valid values:'HTTP1.0','HTTP1.1'.
* `method` - (Optional) Forwarding mode of listener. Valid Values:'RoundRobin', 'LeastConnections'.
* `redirect_alb_listener_id` - (Optional) The ID of the redirect ALB listener.
* `session` - (Optional) session.
* `tls_cipher_policy` - (Optional) TLS cipher policy, valid values:'TlsCipherPolicy1.0','TlsCipherPolicy1.1','TlsCipherPolicy1.2','TlsCipherPolicy1.2-strict','TlsCipherPolicy1.2-moststrict'.

The `session` object supports the following:

* `cookie_name` - (Optional) The name of cookie.
* `cookie_type` - (Optional) The type of cookie, valid values: 'ImplantCookie','RewriteCookie'.
* `session_persistence_period` - (Optional) Session hold timeout. Valid Values:1-86400.
* `session_state` - (Optional) The state of session. Valid Values:'start', 'stop'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `alb_listener_id` - The ID of listener.
* `create_time` - The creation time.


## Import

ALB Listener can be imported using the `id`, e.g.

```
$ terraform import ksyun_alb_listener.example vserver-abcdefg
```

