---
subcategory: "Monitor"
layout: "ksyun"
page_title: "ksyun: ksyun_monitor_alarm_policy"
sidebar_current: "docs-ksyun-resource-monitor_alarm_policy"
description: |-
  Provides a Monitor Alarm Policy resource.
---

# ksyun_monitor_alarm_policy

Provides a Monitor Alarm Policy resource.

#

## Example Usage

```hcl
resource "ksyun_monitor_alarm_policy" "foo" {
  policy_name        = "tf-test-policy"
  product_type       = 0
  policy_type        = 0
  resource_bind_type = 3

  trigger_rules {
    compare       = ">"
    effect_bt     = "00:00"
    effect_et     = "23:59"
    interval      = 5
    item_key      = "cpu.utilizition.total"
    item_name     = "CPU利用率"
    max_count     = 3
    method        = "avg"
    period        = "5m"
    points        = 2
    trigger_value = "90"
    units         = "%"
  }

  instance_ids = []

  user_notice {
    contact_way  = 2
    contact_flag = 2
    contact_id   = 4423
  }

  url_notice = ["https://example.com/webhook"]
}
```

## Argument Reference

The following arguments are supported:

* `policy_name` - (Required, ForceNew) The name of the alarm policy.
* `policy_type` - (Required, ForceNew) Policy type, 0: Normal policy, 1: Default policy.
* `product_type` - (Required, ForceNew) Cloud service category, see [Cloud Service Category](https://docs.ksyun.com/documents/42040).
* `trigger_rules` - (Required, ForceNew) Trigger alarm rules.
* `instance_ids` - (Optional, ForceNew) Instance ID list, required when ResourceBindType=3.
* `project_id` - (Optional, ForceNew) Project group ID, required when ResourceBindType=2.
* `resource_bind_type` - (Optional, ForceNew) Policy associated resource type: 1: All, 2: Project group, 3: Selected instances.
* `url_notice` - (Optional, ForceNew) Alarm callback Webhook addresses, up to 5 can be added.
* `user_notice` - (Optional, ForceNew) Alarm receiving methods.

The `trigger_rules` object supports the following:

* `compare` - (Required) Comparison method, supports >, <, =.
* `interval` - (Required) Alarm interval, unit: minutes, e.g., 5, 10, 30.
* `item_key` - (Required) Monitoring item, e.g., cpu.utilizition.total.
* `item_name` - (Required) Monitoring item name, e.g., CPU utilization.
* `max_count` - (Required) Maximum number of alarm notifications, value range: 1~5.
* `method` - (Required) Statistical method, only supports avg, max, min, sum.
* `period` - (Required) Statistical period, unit: minutes; e.g., 1 minute, 5 minutes, 1 hour correspond to 1m, 5m, 60m respectively.
* `points` - (Required) Consecutive periods.
* `trigger_value` - (Required) Trigger alarm threshold.
* `units` - (Required) Monitoring item unit, e.g., %.
* `effect_bt` - (Optional) Policy effective start time, e.g., 12:00.
* `effect_et` - (Optional) Policy effective end time, e.g., 12:00.

The `user_notice` object supports the following:

* `contact_flag` - (Required) Alarm contact type, 1: Contact group, 2: Contact person.
* `contact_id` - (Required) Contact ID or contact group ID.
* `contact_way` - (Required) Alarm receiving method, 1: Send email, 2: Send SMS, 3: Send email and SMS.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `policy_id` - The ID of the alarm policy.


## Import

Monitor Alarm Policy can be imported using the `policy_id`, e.g.

```
$ terraform import ksyun_monitor_alarm_policy.foo policy_id
```

