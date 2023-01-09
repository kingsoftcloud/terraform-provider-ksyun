---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_host_headers"
sidebar_current: "docs-ksyun-datasource-lb_host_headers"
description: |-
  Provides a list of lb host headers in the current region.
---

# ksyun_lb_host_headers

Provides a list of lb host headers in the current region.

#

## Example Usage

```hcl
data "ksyun_lb_host_headers" "default" {
  output_file = "output_result"
  ids         = []
  listener_id = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of hostheader IDs.
* `listener_id` - (Optional) A list of the listeners.
* `output_file` - (Optional) File name where to save data source results (after running terraform plan).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `host_headers` - An information list of host headers. Each element contains the following attributes:
  * `certificate_id` - The ID of certificate, HTTPS type listener creates this parameter which is not default.
  * `create_time` - The time of creation.
  * `host_header_id` - ID of the host header.
  * `host_header` - the host header.
  * `listener_id` - ID of the listener.
* `total_count` - Total number of host headers that satisfy the condition.


