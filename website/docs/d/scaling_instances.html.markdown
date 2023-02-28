---
subcategory: "Auto Scaling"
layout: "ksyun"
page_title: "ksyun: ksyun_scaling_instances"
sidebar_current: "docs-ksyun-datasource-scaling_instances"
description: |-
  This data source provides a list of ScalingInstance resources in a ScalingGroup.
---

# ksyun_scaling_instances

This data source provides a list of ScalingInstance resources in a ScalingGroup.

#

## Example Usage

```hcl
data "ksyun_scaling_instances" "default" {
  output_file      = "output_result"
  scaling_group_id = "246b37be-5213-49da-a971-xxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required) A scaling group id that the desired ScalingInstance belong to.
* `creation_type` - (Optional) the creation type that desired scalingInstance belong to.
* `health_status` - (Optional) the health status that desired scalingInstance belong to.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `scaling_instance_ids` - (Optional) A list of scaling group ids that the desired ScalingInstance belong to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scaling_instances` - It is a nested type which documented below.
  * `add_time` - The time of creation of ScalingInstance, formatted in RFC3339 time string.
  * `creation_type` - The Creation Type of the desired ScalingInstance.
  * `health_status` - The Health Status of the desired ScalingInstance.
  * `protected_from_detach` - The KEC Instance Protected Model of the desired ScalingInstance.
  * `scaling_instance_id` - The KEC Instance ID of the desired ScalingInstance.
  * `scaling_instance_name` - The KEC Instance Name of the desired ScalingInstance.
* `total_count` - Total number of ScalingInstance resources that satisfy the condition.


