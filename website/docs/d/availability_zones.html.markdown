---
subcategory: "Provider Data Sources"
layout: "ksyun"
page_title: "ksyun: ksyun_availability_zones"
sidebar_current: "docs-ksyun-datasource-availability_zones"
description: |-
  This data source provides a list of available zones in the current region.
---

# ksyun_availability_zones

This data source provides a list of available zones in the current region.

#

## Example Usage

```hcl
data "ksyun_availability_zones" "default" {
  output_file = ""
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `availability_zones` - An information list of AvailabilityZones. Each element contains the following attributes:
  * `availability_zone_name` - Name of the zone.
* `total_count` - Total number of AvailabilityZones that satisfy the condition.


