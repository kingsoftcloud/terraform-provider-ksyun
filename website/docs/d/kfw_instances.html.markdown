---
subcategory: "KFW"
layout: "ksyun"
page_title: "ksyun: ksyun_kfw_instances"
sidebar_current: "docs-ksyun-datasource-kfw_instances"
description: |-
  This data source provides a list of Cloud Firewall Instance resources according to their instance ID, name, and other filters.
---

# ksyun_kfw_instances

This data source provides a list of Cloud Firewall Instance resources according to their instance ID, name, and other filters.

#

## Example Usage

```hcl
data "ksyun_kfw_instances" "default" {
  output_file = "output_result"
  ids         = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Cloud Firewall Instance IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `kfw_instances` - It is a nested type which documented below.
  * `av_status` - AV status (0-stopped, 1-enabled).
  * `bandwidth` - Bandwidth (10-5000M).
  * `cfw_instance_id` - The ID of the Cloud Firewall Instance.
  * `charge_type` - Billing type. Valid values: Monthly (prepaid), Daily (pay-as-you-go, trial).
  * `create_time` - Creation time.
  * `instance_name` - The name of the Cloud Firewall Instance.
  * `instance_type` - Instance type. Valid values: Advanced, Enterprise.
  * `ips_status` - IPS status (0-stopped, 1-enabled).
  * `project_id` - Project ID.
  * `purchase_time` - Purchase duration.
  * `status` - Status (1-creating, 2-running, 3-modifying, 4-stopped, 5-abnormal, 6-unsubscribing).
  * `total_acl_num` - Total number of ACL rules that can be added.
  * `total_eip_num` - Total number of protected IPs.
  * `used_eip_num` - Number of protected IPs in use.
* `total_count` - Total number of Cloud Firewall Instances that satisfy the condition.


