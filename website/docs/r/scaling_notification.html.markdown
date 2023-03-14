---
subcategory: "Auto Scaling"
layout: "ksyun"
page_title: "ksyun: ksyun_scaling_notification"
sidebar_current: "docs-ksyun-resource-scaling_notification"
description: |-
  Provides a ScalingNotification resource.
---

# ksyun_scaling_notification

Provides a ScalingNotification resource.

#

## Example Usage

```hcl
resource "ksyun_scaling_notification" "foo" {
  scaling_group_id           = "541241314798505984"
  scaling_notification_types = ["1", "3"]
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) The ScalingGroup ID of the desired ScalingNotification belong to.
* `scaling_notification_types` - (Required) The List Types of the desired ScalingNotification.Valid Value '1', '2', '3', '4', '5', '6'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `scaling_notification_id` - The ID of the notification.


## Import

```
$ terraform import ksyun_scaling_notification.example scaling-notification-abc123456
```

