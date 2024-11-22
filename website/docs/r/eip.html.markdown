---
subcategory: "EIP"
layout: "ksyun"
page_title: "ksyun: ksyun_eip"
sidebar_current: "docs-ksyun-resource-eip"
description: |-
  Provides an Elastic IP resource.
---

# ksyun_eip

Provides an Elastic IP resource.

## Example Usage

```hcl
data "ksyun_lines" "default" {
  output_file = "output_result1"
  line_name   = "BGP"
}
resource "ksyun_eip" "default" {
  line_id       = "${data.ksyun_lines.default.lines.0.line_id}"
  band_width    = 1
  charge_type   = "PostPaidByPeak"
  purchase_time = 1
  project_id    = 0
}
```

## Argument Reference

The following arguments are supported:

* `band_width` - (Required) The band width of the public address.
* `charge_type` - (Required, ForceNew) The charge type of the Elastic IP address.Valid Values:'PrePaidByMonth','Monthly','PostPaidByPeak','Peak','PostPaidByDay','Daily','PostPaidByTransfer','TrafficMonthly','DailyPaidByTransfer','HourlySettlement','PostPaidByHour','HourlyInstantSettlement','PostpaidByTime'. 
**Notes:** Charge Type have a upgrade, The above-mentioned parameters, **every**, are **valid**. The changes as following:

| Previous Version | Current Version | Description | 
| -------- | -------- | ----------- | 
| PostPaidByPeak | Peak| Pay-as-you-go (monthly peak) | 
 | PostPaidByDay | Daily | Pay-as-you-go (daily) | 
| PostPaidByTransfer | TrafficMonthly | Pay-as-you-go (monthly traffic) |
 | PrePaidByMonth | Monthly | Monthly package | 
|                | DailyPaidByTransfer | Pay-as-you-go (daily traffic) | 
|                | HourlyInstantSettlement | Pay-as-you-go (hourly instant settlement) | 
|                | PostPaidByHour | Pay-as-you-go (hourly billing, monthly settlement) | 
|                | PostpaidByTime | Settlement by times |.
* `line_id` - (Optional, ForceNew) The id of the line.
* `project_id` - (Optional) The id of the project.
* `purchase_time` - (Optional, ForceNew) Purchase time. If charge_type is Monthly or PrePaidByMonth, this is Required.
* `tags` - (Optional) the tags of the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `allocation_id` - the ID of the EIP.
* `band_width_share_id` - the ID of the BWS which the EIP associated.
* `create_time` - creation time of the EIP.
* `instance_id` - the ID of the EIP.
* `instance_type` - The instance type to bind with the EIP.
* `internet_gateway_id` - InternetGateway ID.
* `is_band_width_share` - BWS EIP.
* `network_interface_id` - NetworkInterface ID.
* `public_ip` - The Elastic IP address.
* `state` - state of the EIP.


## Import

EIP can be imported using the id, e.g.

```
$ terraform import ksyun_eip.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

