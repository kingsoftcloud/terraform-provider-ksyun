---
subcategory: "Auto Scaling"
layout: "ksyun"
page_title: "ksyun: ksyun_scaling_group"
sidebar_current: "docs-ksyun-resource-scaling_group"
description: |-
  Provides a ScalingGroup resource.
---

# ksyun_scaling_group

Provides a ScalingGroup resource.

#

## Example Usage

```hcl
resource "ksyun_scaling_group" "foo" {
  subnet_id_set            = [ksyun_subnet.foo.id]
  security_group_id        = ksyun_security_group.foo.id
  scaling_configuration_id = ksyun_scaling_configuration.foo.id
  min_size                 = 0
  max_size                 = 2
  desired_capacity         = 0
  status                   = "Active"
  slb_config_set {
  slb_id = ksyun_lb.foo.id }
  listener_id     = ksyun_lb_listener.foo.id
  server_port_set = [80]
}
}
```

## Argument Reference

The following arguments are supported:

* `desired_capacity` - (Required) The Desire Capacity KEC instance count of the desired ScalingGroup set to.Valid Value 0-1000.
* `max_size` - (Required) The Max KEC instance size of the desired ScalingGroup set to.Valid Value 0-1000.
* `min_size` - (Required) The Min KEC instance size of the desired ScalingGroup set to.Valid Value 0-1000.
* `scaling_configuration_id` - (Required) The Scaling Configuration ID of the desired ScalingGroup set to.
* `remove_policy` - (Optional) The KEC instance remove policy of the desired ScalingGroup set to.Valid Values:'RemoveOldestInstance', 'RemoveNewestInstance'.
* `scaling_group_name` - (Optional) The Name of the desired ScalingGroup.
* `security_group_id_set` - (Optional) The Security Group ID List of the desired ScalingGroup set to.
* `security_group_id` - (Optional) The Security Group ID of the desired ScalingGroup set to.
* `slb_config_set` - (Optional) A list of slb configs.
* `status` - (Optional) The Status of the desired ScalingGroup.Valid Values:'Active', 'UnActive'.
* `subnet_id_set` - (Optional) The Subnet ID Set of the desired ScalingGroup set to.
* `subnet_strategy` - (Optional) The Subnet Strategy of the desired ScalingGroup set to.Valid Values:'balanced-distribution', 'choice-first'.

The `slb_config_set` object supports the following:

* `listener_id` - (Required) The Listener ID of the desired ScalingGroup set to.
* `slb_id` - (Required) The SLB ID of the desired ScalingGroup set to.
* `health_check_type` - (Optional) Health check type, valid values:'slb','kec'.
* `server_port_set` - (Optional) The Server Port Set of the desired ScalingGroup set to.Valid Values 1-65535.
* `weight` - (Optional) The weight of the desired ScalingGroup set to.Valid Values 1-100.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The creation time.
* `instance_num` - The KEC instance Number of the desired ScalingGroup set to.Valid Value 0-10.
* `scaling_configuration_name` - The Scaling Configuration Name of the desired ScalingGroup set to.
* `vpc_id` - The ID of the VPC.


## Import

scalingGroup can be imported using the `id`, e.g.

```
$ terraform import ksyun_scaling_group.example scaling-group-abc123456
```

