---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lbs"
sidebar_current: "docs-ksyun-datasource-lbs"
description: |-
  This data source provides a list of Load Balancer resources according to their Load Balancer ID, VPC ID and Subnet ID.
---

# ksyun_lbs

This data source provides a list of Load Balancer resources according to their Load Balancer ID, VPC ID and Subnet ID.

#

## Example Usage

```hcl
data "ksyun_lbs" "default" {
  output_file = "output_result"
  name_regex  = ""
  ids         = []
  state       = ""
  vpc_id      = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Load Balancer IDs, all the LBs belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter resulting lbs by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `project_id` - (Optional) ID of the project.
* `state` - (Optional) state of the LB.
* `vpc_id` - (Optional) The ID of the VPC linked to the Load Balancers.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `lbs` - It is a nested type which documented below.
  * `access_logs_enabled` - whether accessLogs is enabled or not.
  * `access_logs_s3_bucket` - Bucket for storing access logs.
  * `create_time` - The time of creation.
  * `ip_version` - IP version.
  * `is_waf` - whether it is a waf LB or not.
  * `lb_status` - status of the LB.
  * `lb_type` - Type of the LB.
  * `listeners_count` - ID of the listeners.
  * `load_balancer_id` - ID of the Load Balancer.
  * `load_balancer_name` - Name of the Load Balancer.
  * `load_balancer_state` - start or stop.
  * `project_id` - ID of the project.
  * `public_ip` - public ip address.
  * `state` - associate or disassociate.
  * `subnet_id` - ID of the subnet.
  * `type` - Type of the Load Balancer.
  * `vpc_id` - The ID of the VPC linked to the Load Balancers.
* `total_count` - Total number of Load Balancers that satisfy the condition.


