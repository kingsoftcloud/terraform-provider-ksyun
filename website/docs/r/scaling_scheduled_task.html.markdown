---
subcategory: "Auto Scaling"
layout: "ksyun"
page_title: "ksyun: ksyun_scaling_scheduled_task"
sidebar_current: "docs-ksyun-resource-scaling_scheduled_task"
description: |-
  Provides a ScalingScheduledTask resource.
---

# ksyun_scaling_scheduled_task

Provides a ScalingScheduledTask resource.

#

## Example Usage

```hcl
resource "ksyun_scaling_scheduled_task" "foo" {
  scaling_group_id = "541241314798505984"
  start_time       = "2021-05-01 12:00:00"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required) The ScalingGroup ID of the desired ScalingScheduledTask belong to.
* `start_time` - (Required) The Start Time of the desired ScalingScheduledTask.
* `end_time` - (Optional) The End Time Operator of the desired ScalingScheduledTask.
* `readjust_expect_size` - (Optional) The Readjust Expect Size of the desired ScalingScheduledTask.
* `readjust_max_size` - (Optional) The Readjust Max Size of the desired ScalingScheduledTask.
* `readjust_min_size` - (Optional) The Readjust Min Size of the desired ScalingScheduledTask.
* `recurrence` - (Optional) The Recurrence of the desired ScalingScheduledTask.
* `repeat_cycle` - (Optional) The Repeat Cycle the desired ScalingScheduledTask.
* `repeat_unit` - (Optional) The Repeat Unit of the desired ScalingScheduledTask.
* `scaling_scheduled_task_name` - (Optional) The Name of the desired ScalingScheduledTask.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The creation time.
* `description` - The description of the task.
* `scaling_scheduled_task_id` - The ID of the task.


## Import

```
$ terraform import ksyun_scaling_scheduled_task.example scaling-scheduled-task-abc123456
```

