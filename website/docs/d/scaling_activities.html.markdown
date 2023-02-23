---
subcategory: "Auto Scaling"
layout: "ksyun"
page_title: "ksyun: ksyun_scaling_activities"
sidebar_current: "docs-ksyun-datasource-scaling_activities"
description: |-
  This data source provides a list of ScalingActivity resources in a ScalingGroup.
---

# ksyun_scaling_activities

This data source provides a list of ScalingActivity resources in a ScalingGroup.

#

## Example Usage

```hcl
data "ksyun_scaling_activities" "default" {
  output_file      = "output_result"
  scaling_group_id = "541241314798505984"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required) A ScalingGroup ID that the desired ScalingActivity belong to.
* `end_time` - (Optional) The End Time that the desired ScalingActivity set to.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `start_time` - (Optional) The Start Time that the desired ScalingActivity set to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scaling_activities` - It is a nested type which documented below.
  * `cause` - The cause of the desired ScalingActivity.
  * `description` - The description of the desired ScalingActivity.
  * `end_time` - The end time of the desired ScalingActivity.
  * `error_code` - The error code of the desired ScalingActivity.
  * `scaling_activity_id` - The ID of the desired ScalingActivity.
  * `start_time` - The start time the desired ScalingActivity.
  * `status` - The status of the desired ScalingActivity.
  * `success_instance_list` - The success KEC Instance ID List of the desired ScalingActivity.
  * `type` - The type the desired ScalingActivity.
* `total_count` - Total number of ScalingActivity resources that satisfy the condition.


