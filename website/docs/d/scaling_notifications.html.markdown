---
subcategory: "Auto Scaling"
layout: "ksyun"
page_title: "ksyun: ksyun_scaling_notifications"
sidebar_current: "docs-ksyun-datasource-scaling_notifications"
description: |-
  This data source provides a list of ScalingNotification resources in a ScalingGroup.
---

# ksyun_scaling_notifications

This data source provides a list of ScalingNotification resources in a ScalingGroup.

#

## Example Usage

```hcl
data "ksyun_scaling_notifications" "default" {
  output_file      = "output_result"
  scaling_group_id = "541241314798xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required) A scaling group id that the desired ScalingNotification belong to.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scaling_notifications` - It is a nested type which documented below.
  * `scaling_group_id` - The ScalingGroup ID of the desired ScalingNotification belong to.
  * `scaling_notification_id` - The ID of the ScalingNotification.
  * `scaling_notification_types` - The List Types of the desired ScalingNotification.
* `total_count` - Total number of ScalingNotification resources that satisfy the condition.


