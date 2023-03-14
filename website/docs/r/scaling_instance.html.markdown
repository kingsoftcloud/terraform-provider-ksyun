---
subcategory: "Auto Scaling"
layout: "ksyun"
page_title: "ksyun: ksyun_scaling_instance"
sidebar_current: "docs-ksyun-resource-scaling_instance"
description: |-
  Provides a ScalingInstance resource.
---

# ksyun_scaling_instance

Provides a ScalingInstance resource.

#

## Example Usage

```hcl
resource "ksyun_scaling_instance" "foo" {
  scaling_group_id      = "541241314798505984"
  scaling_instance_id   = "a4ef95c5-e8f1-43f8-912a-758f15064063"
  protected_from_detach = 1
}
```

## Argument Reference

The following arguments are supported:

* `protected_from_detach` - (Required) The KEC Instance Name of the desired ScalingInstance.Valid Value 0,1.
* `scaling_group_id` - (Required, ForceNew) The ScalingGroup ID of the desired ScalingInstance belong to.
* `scaling_instance_id` - (Required, ForceNew) The KEC Instance ID of the desired ScalingInstance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `add_time` - The time of creation of ScalingInstance, formatted in RFC3339 time string.
* `creation_type` - Creation type.
* `health_status` - Health status.
* `scaling_instance_name` - The name of the instance.


## Import

```
$ terraform import ksyun_scaling_instance.example scaling-instance-abc123456
```

