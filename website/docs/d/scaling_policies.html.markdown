---
subcategory: "Auto Scaling"
layout: "ksyun"
page_title: "ksyun: ksyun_scaling_policies"
sidebar_current: "docs-ksyun-datasource-scaling_policies"
description: |-
  This data source provides a list of ScalingPolicy resources in a ScalingGroup.
---

# ksyun_scaling_policies

This data source provides a list of ScalingPolicy resources in a ScalingGroup.

#

## Example Usage

```hcl
data "ksyun_scaling_policies" "default" {
  output_file      = "output_result"
  scaling_group_id = "541241314798xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required) A scaling group id that the desired ScalingPolicy belong to.
* `ids` - (Optional) A list of policy IDs.
* `name_regex` - (Optional) A regex string to filter results by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `scaling_policies_name` - (Optional) The Name that the desired ScalingPolicy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scaling_policies` - It is a nested type which documented below.
  * `adjustment_type` - The Adjustment Type of the desired ScalingPolicy.
  * `adjustment_value` - The Adjustment Value of the desired ScalingPolicy.
  * `comparison_operator` - The Comparison Operator of the desired ScalingPolicy.
  * `cool_down` - The Cool Down of the desired ScalingPolicy.
  * `create_time` - The time of creation of ScalingPolicy, formatted in RFC3339 time string.
  * `dimension_name` - The Dimension Name of the desired ScalingPolicy.
  * `function` - The Function Model of the desired ScalingPolicy.
  * `period` - The Period of the desired ScalingPolicy.
  * `repeat_times` - The Repeat Times of the desired ScalingPolicy.
  * `scaling_group_id` - The ScalingGroup ID of the desired ScalingPolicy belong to.
  * `scaling_policy_id` - The ID of the desired ScalingPolicy.
  * `scaling_policy_name` - The Name of the desired ScalingPolicy.
  * `threshold` - The Threshold of the desired ScalingPolicy.
* `total_count` - Total number of ScalingPolicy resources that satisfy the condition.


