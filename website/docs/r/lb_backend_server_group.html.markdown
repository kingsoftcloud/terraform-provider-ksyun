---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_backend_server_group"
sidebar_current: "docs-ksyun-resource-lb_backend_server_group"
description: |-
  Provides a lb backend server group resource.
---

# ksyun_lb_backend_server_group

Provides a lb backend server group resource.

#

## Example Usage

```hcl
resource "ksyun_lb_backend_server_group" "default" {
  backend_server_group_name = "xuan-tf"
  vpc_id                    = ""
  backend_server_group_type = ""
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) ID of the VPC.
* `backend_server_group_name` - (Optional) The name of backend server group. Default: 'backend_server_group'.
* `backend_server_group_type` - (Optional, ForceNew) The type of backend server group. Valid values: 'Server', 'Mirror'. Default is 'Server'.
* `health_check` - (Optional) Health check information, only the mirror server has this parameter.
* `protocol` - (Optional, ForceNew) The protocol of the backend server group. Valid values: 'TCP', 'UDP', 'HTTP'. Default `HTTP`.

The `health_check` object supports the following:

* `health_check_connect_port` - (Optional) The port of connecting for health check.
* `health_check_state` - (Optional) Status maintained by health examination.Valid Values:'start', 'stop'.
* `healthy_threshold` - (Optional) Health threshold.Valid Values:1-10. Default is 5.
* `host_name` - (Optional) hostname of the health check.
* `interval` - (Optional) Interval of health examination.Valid Values:1-3600. Default is 5.
* `timeout` - (Optional) Health check timeout.Valid Values:1-3600. Default is 4.
* `unhealthy_threshold` - (Optional) Unhealthy threshold.Valid Values:1-10. Default is 4.
* `url_path` - (Optional) Link to HTTP type listener health check.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `backend_server_group_id` - ID of the backend server group.
* `backend_server_number` - number of backend servers.
* `create_time` - creation time of the backend server group.


## Import

LB backend server group can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb_backend_server_group.example fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```

