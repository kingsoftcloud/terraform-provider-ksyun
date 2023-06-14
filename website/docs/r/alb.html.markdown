---
subcategory: "ALB"
layout: "ksyun"
page_title: "ksyun: ksyun_alb"
sidebar_current: "docs-ksyun-resource-alb"
description: |-
  Provides a ALB resource.
---

# ksyun_alb

Provides a ALB resource.

#

## Example Usage

```hcl
resource "ksyun_alb" "default" {
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `alb_name` - (Optional) The name of the ALB.
* `alb_type` - (Optional, ForceNew) The type of the ALB, valid values:'public', 'internal'. Default is 'public'.
* `alb_version` - (Optional, ForceNew) The version of the ALB. valid values:'standard', 'advanced'. Default is 'standard'.
* `charge_type` - (Optional, ForceNew) The charge type, valid values: 'PrePaidByHourUsage'.
* `ip_version` - (Optional, ForceNew) IP version, 'ipv4' or 'ipv6'. Default is 'ipv4'.
* `project_id` - (Optional) The ID of the project.
* `state` - (Optional) The state of the ALB, valid values:'start', 'stop'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The creation time.
* `public_ip` - The public IP address.
* `status` - The status of the ALB.


## Import

BWS can be imported using the id, e.g.

```
$ terraform import ksyun_alb.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

