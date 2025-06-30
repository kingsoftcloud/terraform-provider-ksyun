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

* `alb_type` - (Required, ForceNew) The type of the ALB, valid values:'public', 'internal'.
* `alb_version` - (Required, ForceNew) The version of the ALB. valid values:'standard', 'medium', 'advanced'.
* `charge_type` - (Required, ForceNew) The charge type, valid values: 'PrePaidByHourUsage'.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `alb_name` - (Optional) The name of the ALB.
* `delete_protection` - (Optional) Whether delete protection is enabled or not. Values: `off` or `on`.
* `enable_hpa` - (Optional) Enable hpa.
* `enabled_log` - (Optional) Whether log is enabled or not. Specific `klog_info` field when `enabled_log` is true.
* `enabled_quic` - (Optional, ForceNew) Enable quic.
* `ip_version` - (Optional, ForceNew) IP version, 'ipv4' or 'ipv6'.
* `klog_info` - (Optional) Indicate klog info, including log-project-name and log-pool-name, that use to bind log service for this alb process.
* `modification_protection` - (Optional) Whether modification protection is enabled or not. Values: `off` or `on`.
* `private_ip_address` - (Optional, ForceNew) The private ip address. It not be empty, when 'alb_type' as '**internal**'.
* `project_id` - (Optional) The ID of the project.
* `protocol_layers` - (Optional) The protocol layers of the ALB, valid values: 'L4', 'L7', 'L4-L7'.
* `state` - (Optional) The state of the ALB, Valid Values:'start', 'stop'.
* `status` - (Optional) The status of the ALB.
* `subnet_id` - (Optional, ForceNew) The Id of Subnet that's type is **Reserve**. It not be empty, when 'alb_type' as '**internal**'.
* `tags` - (Optional) the tags of the resource.

The `klog_info` object supports the following:

* `log_pool_name` - (Optional) log pool name.
* `project_name` - (Optional) log project name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The creation time.
* `public_ip` - The public IP address.


## Import

`ksyun_alb` can be imported using the id, e.g.

```
$ terraform import ksyun_alb.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

