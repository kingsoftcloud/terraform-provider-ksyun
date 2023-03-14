---
subcategory: "Auto Scaling"
layout: "ksyun"
page_title: "ksyun: ksyun_scaling_policy"
sidebar_current: "docs-ksyun-resource-scaling_policy"
description: |-
  Provides a ScalingPolicy resource.
---

# ksyun_scaling_policy

Provides a ScalingPolicy resource.

#

## Example Usage

```hcl
resource "ksyun_scaling_policy" "foo" {
  scaling_group_id = "541241314798505984"
  threshold        = 20
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) The ScalingGroup ID of the desired ScalingNotification belong to.
* `adjustment_type` - (Optional) The Adjustment Type of the desired ScalingPolicy.Valid Value 'TotalCapacity', 'QuantityChangeInCapacity', 'PercentChangeInCapacity'.
* `adjustment_value` - (Optional) The Adjustment Value of the desired ScalingPolicy.Valid Value -100 ~ 100.
* `comparison_operator` - (Optional) The Comparison Operator of the desired ScalingPolicy.Valid Value 'Greater', 'EqualOrGreater', 'Less', 'EqualOrLess', 'Equal', 'NotEqual'.
* `cool_down` - (Optional) The Cool Down of the desired ScalingPolicy.Min is 60.
* `dimension_name` - (Optional) The Dimension Name of the desired ScalingPolicy.Valid Value 'cpu_usage', 'mem_usage', 'net_outtraffic', 'net_intraffic', 'listener_outtraffic', 'listener_intraffic'.
* `function` - (Optional) The Function Model of the desired ScalingPolicy.Valid Value 'avg', 'min', 'max'.
* `period` - (Optional) The Period of the desired ScalingPolicy.Min is 60.
* `repeat_times` - (Optional) The Repeat Times of the desired ScalingPolicy.Valid Value 1-10.
* `scaling_policy_name` - (Optional) The Name of the desired ScalingPolicy.
* `threshold` - (Optional) The Threshold of the desired ScalingPolicy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The creation time.
* `scaling_policy_id` - The ID of the scaling policy.


## Import

```
$ terraform import ksyun_scaling_policy.example scaling-policy-abc123456
```

