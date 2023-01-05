---
subcategory: "Bare Metal"
layout: "ksyun"
page_title: "ksyun: ksyun_bare_metal_raid_attributes"
sidebar_current: "docs-ksyun-datasource-bare_metal_raid_attributes"
description: |-
  This data source provides a list of Bare Metal Raid Attributes resources according to their Bare Metal Raid Attribute ID.
---

# ksyun_bare_metal_raid_attributes

This data source provides a list of Bare Metal Raid Attributes resources according to their Bare Metal Raid Attribute ID.

#

## Example Usage

```hcl
# Get  bare metal_raid_attributes

data "ksyun_bare_metal_raid_attributes" "default" {
  output_file = "output_result"
}
```

## Argument Reference

The following arguments are supported:

* `host_type` - (Optional) A list of Bare Metal Raid Attribute Host Types.
* `name_regex` - (Optional) A regex string to filter results by name of Bare Metal Raid template.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `raid_attributes` - Total number of Bare Metal Raid Attributes that satisfy the condition.
  * `create_time` - The time of creation for Bare Metal Raid template.
  * `disk_set` - list of disks that used raid template.
    * `disk_attribute` - attribute of the disk.
    * `disk_count` - count of disks.
    * `disk_id` - ID of the disk.
    * `disk_space` - space of the data disk.
    * `disk_type` - type of the disk.
    * `raid` - raid level.
    * `space` - available Space.
    * `system_disk_space` - space of the system disk.
  * `host_type` - host type of the Bare Metal.
  * `raid_id` - ID of the raid.
  * `template_name` - name of the Bare Metal Raid template.
* `total_count` - Total number of Bare Metal Raid Attributes that satisfy the condition.


