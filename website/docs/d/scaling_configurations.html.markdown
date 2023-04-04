---
subcategory: "Auto Scaling"
layout: "ksyun"
page_title: "ksyun: ksyun_scaling_configurations"
sidebar_current: "docs-ksyun-datasource-scaling_configurations"
description: |-
  This data source provides a list of ScalingConfiguration resources.
---

# ksyun_scaling_configurations

This data source provides a list of ScalingConfiguration resources.

#

## Example Usage

```hcl
data "ksyun_scaling_configurations" "default" {
  output_file                = "output_result"
  ids                        = []
  project_ids                = []
  scaling_configuration_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ScalingConfiguration IDs, all the ScalingConfiguration resources belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter results by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `project_ids` - (Optional) A list of Project id that the desired ScalingConfiguration belongs to.
* `scaling_configuration_name` - (Optional) The Name of ScalingConfiguration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scaling_configurations` - It is a nested type which documented below.
  * `address_band_width` - IP band width.
  * `address_project_id` - The project ID of the IP address.
  * `band_width_share_id` - The ID of the BWS.
  * `charge_type` - charge type.
  * `cpu` - The CPU core size of the desired ScalingConfiguration.
  * `create_time` - The time of creation of ScalingGroup, formatted in RFC3339 time string.
  * `data_disk_gb` - The Local Volume GB size of the desired ScalingConfiguration.
  * `data_disks` - It is a nested type which documented below.
    * `delete_with_instance` - The Flag with delete EBS Data Disk when KEC Instance destroy.
    * `disk_size` - The EBS Data Disk Size of the desired data_disk.
    * `disk_type` - The EBS Data Disk Type of the desired data_disk.
  * `gpu` - The GPU core size the desired ScalingConfiguration.
  * `image_id` - The System Image Id of the desired ScalingConfiguration.
  * `instance_name_suffix` - The kec instance name suffix of the desired ScalingConfiguration.
  * `instance_name_time_suffix` - The kec instance name suffix of the desired ScalingConfiguration.
  * `instance_name` - The KEC instance name of the desired ScalingConfiguration.
  * `instance_type` - The KEC instance type of the desired ScalingConfiguration.
  * `keep_image_login` - The Flag with image login set of the desired ScalingConfiguration.
  * `key_id` - The SSH key set of the desired ScalingConfiguration.
  * `line_id` - The ID of the line.
  * `mem` - The Memory GB size of the desired ScalingConfiguration.
  * `need_monitor_agent` - The Monitor agent flag desired ScalingConfiguration.
  * `need_security_agent` - The Security agent flag desired ScalingConfiguration.
  * `project_id` - The Project Id of the desired ScalingConfiguration belong to.
  * `scaling_configuration_id` - The ID of ScalingConfiguration.
  * `scaling_configuration_name` - The Name of the desired ScalingConfiguration.
  * `system_disk_size` - System disk size.
  * `system_disk_type` - System disk type.
  * `user_data` - The user data of the desired ScalingConfiguration.
* `total_count` - Total number of ScalingConfiguration resources that satisfy the condition.


