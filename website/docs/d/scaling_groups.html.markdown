---
subcategory: "Auto Scaling"
layout: "ksyun"
page_title: "ksyun: ksyun_scaling_groups"
sidebar_current: "docs-ksyun-datasource-scaling_groups"
description: |-
  This data source provides a list of ScalingGroup resources .
---

# ksyun_scaling_groups

This data source provides a list of ScalingGroup resources .

#

## Example Usage

```hcl
data "ksyun_scaling_groups" "default" {
  output_file = "output_result"
  vpc_id      = "246b37be-5213-49da-a971-xxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ScalingGroup IDs, all the ScalingGroup resources belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter results by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `scaling_configuration_id` - (Optional) The Scaling Configuration ID of the desired ScalingGroup set to.
* `scaling_group_name` - (Optional) The Name of the desired ScalingGroup.
* `vpc_id` - (Optional) The VPC ID of the desired ScalingGroup set to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scaling_groups` - It is a nested type which documented below.
  * `create_time` - The time of creation of ScalingGroup, formatted in RFC3339 time string.
  * `desired_capacity` - The Desire Capacity KEC instance count of the desired ScalingGroup set to.
  * `instance_num` - The KEC instance Number of the desired ScalingGroup set to.
  * `max_size` - The Min KEC instance size of the desired ScalingGroup set to.
  * `min_size` - The Min KEC instance size of the desired ScalingGroup set to.
  * `remove_policy` - The KEC instance remove policy of the desired ScalingGroup set to.
  * `scaling_configuration_id` - The Scaling Configuration ID of the desired ScalingGroup set to.
  * `scaling_configuration_name` - The Scaling Configuration Name of the desired ScalingGroup set to.
  * `scaling_group_id` - The Security Group ID of the desired ScalingGroup set to.
  * `scaling_group_name` - The Name of the desired ScalingGroup.
  * `security_group_id_set` - A list of the Security Group IDs.
  * `security_group_id` - The Security Group ID of the desired ScalingGroup set to.
  * `slb_config_set` - The SLB Config Set of the desired ScalingGroup set to.
    * `health_check_type` - The health check type of the desired ScalingGroup set to.
    * `listener_id` - The Listener ID of the desired ScalingGroup set to.
    * `server_port_set` - The Server Port Set of the desired ScalingGroup set to.
    * `slb_id` - The SLB ID of the desired ScalingGroup set to.
    * `weight` - The weight of the desired ScalingGroup set to.
  * `status` - The Status of the desired ScalingGroup.
  * `subnet_id_set` - The Subnet ID Set of the desired ScalingGroup set to.
  * `subnet_strategy` - The Subnet Strategy of the desired ScalingGroup set to.
  * `vpc_id` - The VPC ID of the desired ScalingGroup set to.
* `total_count` - Total number of ScalingGroup resources that satisfy the condition.


