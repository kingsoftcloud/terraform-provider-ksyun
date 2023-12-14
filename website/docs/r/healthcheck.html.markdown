---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_healthcheck"
sidebar_current: "docs-ksyun-resource-healthcheck"
description: |-
  Provides an Health Check resource.
---

# ksyun_healthcheck

Provides an Health Check resource.

#

## Example Usage

```hcl
resource "ksyun_healthcheck" "default" {
  listener_id          = "537e2e7b-0007-4a75-9749-882167dbc93d"
  health_check_state   = "stop"
  healthy_threshold    = 2
  interval             = 20
  timeout              = 200
  unhealthy_threshold  = 2
  url_path             = "/monitor"
  is_default_host_name = true
  host_name            = "www.ksyun.com"
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, ForceNew) The id of the listener.
* `health_check_connect_port` - (Optional) The port of connecting for health check.
* `health_check_state` - (Optional) Status maintained by health examination.Valid Values:'start', 'stop'.
* `healthy_threshold` - (Optional) Health threshold.Valid Values:1-10. Default is 5.
* `host_name` - (Optional) The service host name of the health check, which is available only for the HTTP or HTTPS health check.
* `http_method` - (Optional) The http requests' method. Valid Value: GET|HEAD.
* `interval` - (Optional) Interval of health examination.Valid Values:1-3600. Default is 5.
* `is_default_host_name` - (Optional) Whether the host name is default or not.
* `lb_type` - (Optional, ForceNew) The type of listener. Valid Value: `Alb` and `Slb`. Default: `Slb`.
* `timeout` - (Optional) Health check timeout.Valid Values:1-3600. Default is 4.
* `unhealthy_threshold` - (Optional) Unhealthy threshold.Valid Values:1-10. Default is 4.
* `url_path` - (Optional) Link to HTTP type listener health check.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `health_check_id` - ID of the health check.
* `listener_protocol` - protocol of the listener.


## Import

HealthCheck can be imported using the id, e.g.

```
$ terraform import ksyun_healthcheck.default ${lb_type}:67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

