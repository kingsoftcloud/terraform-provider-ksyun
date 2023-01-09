---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_instances"
sidebar_current: "docs-ksyun-datasource-instances"
description: |-
  This data source providers a list of instance resources according to their availability zone, instance ID.
---

# ksyun_instances

This data source providers a list of instance resources according to their availability zone, instance ID.

#

## Example Usage

```hcl
# Get  instances

data "ksyun_instances" "default" {
  output_file = "output_result"
  ids         = []
  search      = ""
  project_id  = []
  network_interface {
    network_interface_id = []
    subnet_id            = []
    group_id             = []
  }
  instance_state {
    name = []
  }
  availability_zone {
    name = []
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) the availability zone that the instance locates at.
* `ids` - (Optional) A list of instance IDs.
* `instance_state` - (Optional) The state of instance.
* `name_regex` - (Optional) A regex string to filter results by instance name.
* `network_interface` - (Optional) a list of network interface.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `project_id` - (Optional) One or more project IDs.
* `search` - (Optional) A regex string to filter results by instance name or privateIpAddress.
* `subnet_id` - (Optional) The ID of subnet linked to the instance.
* `vpc_id` - (Optional) The ID of VPC linked to the instance.

The `availability_zone` object supports the following:

* `name` - (Optional) 

The `instance_state` object supports the following:

* `name` - (Optional) name of the state.

The `network_interface` object supports the following:

* `group_id` - (Optional) The ID of security group linked to the network interface.
* `network_interface_id` - (Optional) the ID of the network interface.
* `subnet_id` - (Optional) The ID of VPC linked to the network interface.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - a list of instance.
  * `auto_scaling_type` - type of auto scaling.
  * `availability_zone_name` - Availability zone name.
  * `availability_zone` - Availability zone name.
  * `charge_type` - Instance charge type.
  * `creation_date` - The time of creation for instance.
  * `data_disks` - a list of the data disks.
    * `delete_with_instance` - Decides whether the disk is deleted with instance.
    * `disk_id` - ID of the data disk.
    * `disk_size` - size of the data disk.
    * `disk_type` - type of the data disk.
  * `image_id` - ID of the image.
  * `instance_configure` - the configure of the instance.
    * `data_disk_gb` - size of the data disk.
    * `data_disk_type` - type of the data disk.
    * `g_p_u` - the number of the gpu.
    * `memory_gb` - memory capacity.
    * `root_disk_gb` - size of the root disk.
    * `v_c_p_u` - the number of the vcpu.
  * `instance_count` - count of the instance.
  * `instance_id` - the ID of the instance.
  * `instance_name` - the name of the instance.
  * `instance_state` - state of the instance.
    * `name` - name of the state.
  * `instance_type` - type of the instance.
  * `is_show_sriov_net_support` - whether support networking enhancement.
  * `key_id` - The certificate id of the instance.
  * `monitoring` - state of the monitoring.
    * `state` - name of the state.
  * `network_interface_set` - a list of network interface.
    * `d_n_s1` - The dns1 of the network interface.
    * `d_n_s2` - The dns2 of the network interface.
    * `group_set` - a list of the security group.
      * `group_id` - ID of the security group.
    * `mac_address` - MAC address.
    * `network_interface_id` - ID of the network interface.
    * `network_interface_type` - type of the network interface.
    * `private_ip_address` - private ip address of the network interface.
    * `public_ip` - public ip address of the network interface.
    * `security_group_set` - a list of the security group.
      * `security_group_id` - ID of the security group.
    * `subnet_id` - ID of the subnet.
  * `private_ip_address` - Instance private IP address.
  * `product_type` - product type of the instance.
  * `product_what` - whether the instance is trial or not.
  * `project_id` - The project instance belongs to.
  * `sriov_net_support` - whether support networking enhancement.
  * `stopped_mode` - stopped mode.
  * `subnet_id` - The ID of subnet linked to the instance.
  * `system_disk` - System disk information.
    * `disk_size` - size of the system disk.
    * `disk_type` - type of the system disk.
  * `vpc_id` - The ID of VPC linked to the instance.
* `total_count` - Total number of instance that satisfy the condition.


