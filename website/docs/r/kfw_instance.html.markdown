---
subcategory: "KFW"
layout: "ksyun"
page_title: "ksyun: ksyun_kfw_instance"
sidebar_current: "docs-ksyun-resource-kfw_instance"
description: |-
  Provides a Cloud Firewall Instance resource.
---

# ksyun_kfw_instance

Provides a Cloud Firewall Instance resource.

#

## Example Usage

```hcl
resource "ksyun_kfw_instance" "default" {
  instance_name = "test-kfw-instance"
  instance_type = "Advanced"
  bandwidth     = 50
  total_eip_num = 50
  charge_type   = "Monthly"
  project_id    = "0"
  purchase_time = 1
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) Bandwidth (10-5000M). Must be a multiple of 5M. Advanced: minimum 10M, Enterprise: minimum 50M.
* `charge_type` - (Required, ForceNew) Billing type. Valid values: Monthly (prepaid), Daily (pay-as-you-go, trial).
* `instance_type` - (Required, ForceNew) Instance type. Valid values: Advanced, Enterprise.
* `project_id` - (Required) Project ID. Length 0-36 characters, supports letters, numbers, hyphens(-).
* `total_eip_num` - (Required) Total number of protected IPs. Range: 1-500. Advanced: minimum 20, Enterprise: minimum 50.
* `instance_name` - (Optional) The name of the Cloud Firewall Instance. Length 0-64 characters, supports Chinese, English, numbers.
* `purchase_time` - (Optional, ForceNew) Purchase duration. Required when charge_type is Monthly, range: 1-36 months. Required when charge_type is Daily and ProductWhat is 2, range: 1-14 days.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `av_status` - AV status (0-stopped, 1-enabled).
* `cfw_instance_id` - The ID of the Cloud Firewall Instance.
* `create_time` - Creation time.
* `ips_status` - IPS status (0-stopped, 1-enabled).
* `status` - Status (1-creating, 2-running, 3-modifying, 4-stopped, 5-abnormal, 6-unsubscribing).
* `total_acl_num` - Total number of ACL rules that can be added.
* `used_eip_num` - Number of protected IPs in use.


## Import

Cloud Firewall Instance can be imported using the `cfw_instance_id`, e.g.

```
$ terraform import ksyun_kfw_instance.default cfw_instance_id
```

