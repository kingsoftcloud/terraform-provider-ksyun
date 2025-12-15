---
subcategory: "ALB"
layout: "ksyun"
page_title: "ksyun: ksyun_alb_backend_server_group"
sidebar_current: "docs-ksyun-resource-alb_backend_server_group"
description: |-
  Provides an alb backend server group resource.
---

# ksyun_alb_backend_server_group

Provides an alb backend server group resource.

#

## Example Usage

```hcl
resource "ksyun_vpc" "test" {
  vpc_name   = "tf-alb-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_alb_backend_server_group" "foo" {
  name                = "tf-alb-bsg"
  vpc_id              = ksyun_vpc.test.id
  upstream_keepalive  = "adaptation"
  backend_server_type = "Host"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) ID of the VPC.
* `backend_server_type` - (Optional, ForceNew) The type of backend server. Valid values: 'Host', 'DirectConnect'. Default is 'Host'.
* `health_check` - (Optional) Health check information.
* `method` - (Optional) Forwarding mode of listener. Valid Values:'RoundRobin', 'LeastConnections', 'MasterSlave', 'QUIC_CID', 'IPHash'.
* `name` - (Optional) The name of alb backend server group. Default: 'ksc_bsg'.
* `protocol` - (Optional, ForceNew) The protocol of backend server. Valid values: 'HTTP', 'gRPC', 'TCP', 'UDP', 'HTTPS'. Default is 'HTTP'.
* `session` - (Optional) Whether keeps session. Specific `session` block, if keeps session.
* `upstream_keepalive` - (Optional) The upstream keepalive type. Valid Value: `adaptation`, `keepalive`, `shortconnection`.

The `health_check` object supports the following:

* `health_check_connect_port` - (Optional) The port of connecting for health check.
* `health_check_exp` - (Optional) The expected response of health check.
* `health_check_req` - (Optional) The request of health check.
* `health_check_state` - (Optional) Status maintained by health examination.Valid Values:'start', 'stop'.
* `health_code` - (Optional) The health check code.
* `health_protocol` - (Optional) The health check protocol. Valid values: 'HTTP', 'TCP', 'ICMP', 'UDP'.
* `healthy_threshold` - (Optional) Health threshold.Valid Values:1-10. Default is 5.
* `host_name` - (Optional) hostname of the health check.
* `http_method` - (Optional) The http requests' method. Valid Value: GET|HEAD.
* `interval` - (Optional) Interval of health examination.Valid Values:1-3600. Default is 5.
* `timeout` - (Optional) Health check timeout.Valid Values:1-3600. Default is 4.
* `unhealthy_threshold` - (Optional) Unhealthy threshold.Valid Values:1-10. Default is 4.
* `url_path` - (Optional) Link to HTTP type listener health check.

The `session` object supports the following:

* `cookie_name` - (Optional) The name of cookie.
* `cookie_type` - (Optional) The type of cookie, valid values: 'ImplantCookie','RewriteCookie'.
* `session_persistence_period` - (Optional) Session hold timeout. Valid Values:1-86400.
* `session_state` - (Optional) The state of session. Valid Values:'start', 'stop'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `backend_server_group_id` - ID of the alb backend server group.
* `backend_server_number` - number of backend servers.
* `create_time` - creation time of the alb backend server group.


## Import

ALB backend server group can be imported using the `id`, e.g.

```
$ terraform import ksyun_alb_backend_server_group.default fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```

