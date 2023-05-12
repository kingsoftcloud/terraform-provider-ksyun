---
subcategory: "Auto Scaling"
layout: "ksyun"
page_title: "ksyun: ksyun_scaling_scheduled_tasks"
sidebar_current: "docs-ksyun-datasource-scaling_scheduled_tasks"
description: |-
  This data source provides a list of ScalingScheduledTask resources in a ScalingGroup.
---

# ksyun_scaling_scheduled_tasks

This data source provides a list of ScalingScheduledTask resources in a ScalingGroup.

#

## Example Usage

```hcl
data "ksyun_scaling_scheduled_tasks" "default" {
  output_file      = "output_result"
  scaling_group_id = "541241314798xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required) A scaling group id that the desired ScalingScheduledTask belong to.
* `ids` - (Optional) A list of resource IDs.
* `name_regex` - (Optional) A regex string to filter results by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `scaling_scheduled_task_name` - (Optional) The Name that the desired ScalingScheduledTask.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scaling_scheduled_tasks` - It is a nested type which documented below.
  * `create_time` - The time of creation of ScalingScheduledTask, formatted in RFC3339 time string.
  * `description` - The Description of the desired ScalingScheduledTask.
  * `end_time` - The End Time Operator of the desired ScalingScheduledTask.
  * `readjust_expect_size` - The Readjust Expect Size of the desired ScalingScheduledTask.
  * `readjust_max_size` - The Readjust Max Size of the desired ScalingScheduledTask.
  * `readjust_min_size` - The Readjust Min Size of the desired ScalingScheduledTask.
  * `recurrence` - The Recurrence of the desired ScalingScheduledTask.
  * `repeat_cycle` - The Repeat Cycle the desired ScalingScheduledTask.
  * `repeat_unit` - The Repeat Unit of the desired ScalingScheduledTask.
  * `scaling_group_id` - The ScalingGroup ID of the desired ScalingScheduledTask belong to.
  * `scaling_scheduled_task_id` - The ID of the desired ScalingScheduledTask.
  * `scaling_scheduled_task_name` - The Name of the desired ScalingScheduledTask.
  * `start_time` - The Start Time of the desired ScalingScheduledTask.
* `total_count` - Total number of ScalingScheduledTask resources that satisfy the condition.


