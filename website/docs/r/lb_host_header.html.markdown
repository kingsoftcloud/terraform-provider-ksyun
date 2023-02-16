---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_host_header"
sidebar_current: "docs-ksyun-resource-lb_host_header"
description: |-
  Provides a lb host header resource.
---

# ksyun_lb_host_header

Provides a lb host header resource.

#

## Example Usage

```hcl
resource "ksyun_lb_host_header" "default" {
  listener_id    = "xxxx"
  host_header    = "tf-xuan"
  certificate_id = ""
}
```

EIP can be imported using the id, e.g.

```hcl
terraform import ksyun_lb_host_header.default 67 b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

## Argument Reference

The following arguments are supported:

* `host_header` - (Required, ForceNew) The host header.
* `listener_id` - (Required, ForceNew) The ID of the listener.
* `certificate_id` - (Optional) The ID of the certificate, HTTPS type listener creates this parameter which is not default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time when the host header was created.
* `host_header_id` - The host header id.
* `listener_protocol` - The protocol of the listener.


