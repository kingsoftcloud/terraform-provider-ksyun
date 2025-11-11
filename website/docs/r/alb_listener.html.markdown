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
# network and security group configuration
resource "ksyun_vpc" "example" {
  vpc_name   = "tf-alb-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "example" {
  subnet_name       = "tf-alb-subnet"
  cidr_block        = "10.0.1.0/24"
  subnet_type       = "Normal"
  availability_zone = var.az
  vpc_id            = ksyun_vpc.example.id
}

resource "ksyun_security_group" "example" {
  vpc_id              = ksyun_vpc.example.id
  security_group_name = "tf_sg"
}

resource "ksyun_security_group_entry" "example" {
  security_group_id = ksyun_security_group.example.id
  cidr_block        = "0.0.0.0/0"
  protocol          = "ip"
  direction         = "in"
}

# ---------------------------------------------
# alb backend server group for default forward rule group
resource "ksyun_alb_backend_server_group" "foo" {
  name                = "tf-alb-bsg"
  vpc_id              = ksyun_vpc.example.id
  upstream_keepalive  = "adaptation"
  backend_server_type = "Host"
}

# ---------------------------------------------
# resource ksyun alb
resource "ksyun_alb" "example" {
  alb_name    = "tf-alb-example-1"
  alb_version = "standard"
  alb_type    = "public"
  state       = "start"
  charge_type = "PrePaidByHourUsage"
  vpc_id      = ksyun_vpc.example.id
  project_id  = 0
  enabled_log = false
  ip_version  = "ipv4"
}

# query your certificates on ksyun
data "ksyun_certificates" "listener_cert" {
  name_regex = "test"
}

resource "ksyun_alb_listener" "example" {
  alb_id             = ksyun_alb.example.id
  alb_listener_name  = "alb-example-listener"
  protocol           = "HTTPS"
  port               = 8099
  alb_listener_state = "start"
  certificate_id     = data.ksyun_certificates.listener_cert.certificates.0.certificate_id

  session {
    cookie_type                = "ImplantCookie"
    cookie_name                = "KLBRSIDdad"
    session_state              = "start"
    session_persistence_period = 3100
  }

  # default forward rule setting
  default_forward_rule {
    backend_server_group_id = ksyun_alb_backend_server_group.foo.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `alb_id` - (Required, ForceNew) The ID of the ALB.
* `port` - (Required, ForceNew) The protocol port of listener.
* `protocol` - (Required, ForceNew) The protocol of listener. Valid Values: 'HTTP', 'HTTPS', 'TCP', 'UDP', 'TCPSSL'.
* `alb_listener_name` - (Optional) The name of the listener.
* `alb_listener_state` - (Optional) The state of listener.Valid Values:'start', 'stop'.
* `ca_certificate_id` - (Optional) The ID of Client's CA certificate.
* `ca_enabled` - (Optional) Whether enable to CA certificate.
* `certificate_id` - (Optional) The ID of certificate.
* `config_content` - (Optional) The custom configure for listener. [The details](https://docs.ksyun.com/documents/42615?type=3).
* `default_forward_rule` - (Optional) The default forward rule group.
* `enable_http2` - (Optional) whether enable to HTTP2.
* `enable_quic_upgrade` - (Optional) Whether enable to QUIC upgrade.
* `http_protocol` - (Optional, **Deprecated**) This field will be removed soon. Please use 'enable_http2' instead to choose a protocol. Backend Protocol, valid values:'HTTP1.0','HTTP1.1'.
* `method` - (Optional) Forwarding mode of listener. Valid Values:'RoundRobin', 'LeastConnections'.
* `quic_listener_id` - (Optional) The ID of QUIC listener.
* `redirect_alb_listener_id` - (Optional, **Deprecated**) This parameter is moved to 'default_forward_rule' block. The ID of the redirect ALB listener.
* `server_group_id` - (Optional) The ID of the backend server group.
* `session` - (Optional) Whether keeps session. Specific `session` block, if keeps session.
* `tls_cipher_policy` - (Optional) TLS cipher policy, valid values:'TlsCipherPolicy1.0','TlsCipherPolicy1.1','TlsCipherPolicy1.2','TlsCipherPolicy1.2-strict','TlsCipherPolicy1.2-most-strict-with1.3'.

The `default_forward_rule` object supports the following:

* `backend_server_group_id` - (Optional) The backend server group id for default forward rule group.
* `fixed_response_config` - (Optional) 
* `redirect_alb_listener_id` - (Optional) The ID of the alternative redirect ALB listener.
* `redirect_http_code` - (Optional) The http code for redirect action. Valid Values: 301|302|307.
* `rewrite_config` - (Optional) The config of rewrite.
* `type` - (Optional, ForceNew) The type of default forward rule group. Valid Values: 'Redirect', 'FixedResponse', 'Rewrite', 'ForwardGroup.

The `rewrite_config` object supports the following:

* `http_host` - (Optional) The host of the rewrite.
* `query_string` - (Optional) The query string of the rewrite.
* `url` - (Optional) The url of the rewrite.

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
* `redirect_listener_name` - The redirect listener name.


## Import

ALB Listener can be imported using the `id`, e.g.

```
$ terraform import ksyun_alb_listener.example vserver-abcdefg
```

