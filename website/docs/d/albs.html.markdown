---
subcategory: "ALB"
layout: "ksyun"
page_title: "ksyun: ksyun_albs"
sidebar_current: "docs-ksyun-datasource-albs"
description: |-
  This data source provides a list of ALB resources according to their ALB ID.
---

# ksyun_albs

This data source provides a list of ALB resources according to their ALB ID.

#

## Example Usage

```hcl
data "ksyun_albs" "default" {
  output_file = "output_result"
  ids         = []
  vpc_id      = []
  state       = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ALB IDs, all the ALBs belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `state` - (Optional) One or more state.
* `vpc_id` - (Optional) One or more VPC IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `albs` - An information list of ALB. Each element contains the following attributes:
  * `alb_id` - The ID of the ALB.
  * `alb_name` - The name of the ALB.
  * `alb_type` - The type of the ALB.
  * `alb_version` - The version of the ALB.
  * `charge_type` - The charge type.
  * `create_time` - The creation time.
  * `enabled_log` - whether log is enabled or not.
  * `id` - ID of the ALB.
  * `ip_version` - IP version, 'ipv4' or 'ipv6'.
  * `klog_info` - klog info.
    * `account_id` - account id.
    * `log_pool_name` - log pool name.
    * `project_name` - log project name.
  * `project_id` - The ID of the project.
  * `public_ip` - The public IP address.
  * `state` - The state of the ALB.
  * `status` - The status of the ALB.
  * `vpc_id` - The ID of the VPC.
* `total_count` - Total number of ALBs that satisfy the condition.


