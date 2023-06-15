---
subcategory: "ALB"
layout: "ksyun"
page_title: "ksyun: ksyun_alb_listener_cert_group"
sidebar_current: "docs-ksyun-resource-alb_listener_cert_group"
description: |-
  Provides a ALB Listener cert group resource.
---

# ksyun_alb_listener_cert_group

Provides a ALB Listener cert group resource.

#

## Example Usage

```hcl
resource "ksyun_alb_listener_cert_group" "default" {
}
```

## Argument Reference

The following arguments are supported:

* `alb_listener_id` - (Required, ForceNew) The ID of the ALB Listener.
* `alb_listener_cert_set` - (Optional) certificate list.

The `alb_listener_cert_set` object supports the following:

* `certificate_id` - (Required) The ID of the certificate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `alb_listener_cert_group_id` - The ID of the ALB Listener Cert Group.


## Import

ALB Listener Cert Group can be imported using the `id`, e.g.

```
$ terraform import ksyun_alb_listener_cert_group.example vserver-abcdefg
```

