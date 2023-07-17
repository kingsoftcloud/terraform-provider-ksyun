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
resource "ksyun_vpc" "default" {
  vpc_name   = "tf_alb_test_vpc"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_alb" "default" {
  alb_name    = "tf_test_alb1"
  alb_version = "standard"
  alb_type    = "public"
  state       = "start"
  charge_type = "PrePaidByHourUsage"
  vpc_id      = ksyun_vpc.default.id
  project_id  = 0
}

data "ksyun_lines" "default" {
  output_file = "output_result1"
  line_name   = "BGP"
}

resource "ksyun_eip" "foo" {
  line_id       = data.ksyun_lines.default.lines.0.line_id
  band_width    = 1
  charge_type   = "PostPaidByDay"
  purchase_time = 1
  project_id    = 0
}

resource "ksyun_eip_associate" "eip_bind" {
  allocation_id = ksyun_eip.foo.id
  instance_id   = ksyun_alb.foo.id
  instance_type = "Slb"
}
```

## Argument Reference

The following arguments are supported:

* `alb_type` - (Required, ForceNew) The type of the ALB, valid values:'public', 'internal''.
* `alb_version` - (Required, ForceNew) The version of the ALB. valid values:'standard', 'advanced'.
* `charge_type` - (Required, ForceNew) The charge type, valid values: 'PrePaidByHourUsage'.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `alb_name` - (Optional) The name of the ALB.
* `enabled_log` - (Optional) whether log is enabled or not.
* `ip_version` - (Optional, ForceNew) IP version, 'ipv4' or 'ipv6'.
* `klog_info` - (Optional) klog info.
* `project_id` - (Optional) The ID of the project.
* `state` - (Optional) The state of the ALB, valid values:'start', 'stop'.

The `klog_info` object supports the following:

* `log_pool_name` - (Optional) log pool name.
* `project_name` - (Optional) log project name.

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

