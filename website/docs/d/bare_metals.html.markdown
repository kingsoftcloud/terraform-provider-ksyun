---
subcategory: "Bare Metal"
layout: "ksyun"
page_title: "ksyun: ksyun_bare_metals"
sidebar_current: "docs-ksyun-datasource-bare_metals"
description: |-
  This data source provides a list of Bare Metal resources according to their Bare Metal ID.
---

# ksyun_bare_metals

This data source provides a list of Bare Metal resources according to their Bare Metal ID.

#

## Example Usage

```hcl
# Get  bare metals

data "ksyun_bare_metals" "default" {
  output_file     = "output_result"
  ids             = []
  vpc_id          = ["bfec0f43-9e5a-4f06-b7a1-df4768c1cd6f"]
  project_id      = []
  host_name       = []
  subnet_id       = []
  cabinet_id      = []
  epc_host_status = []
  os_name         = []
  product_type    = []
}
```

## Argument Reference

The following arguments are supported:

* `cabinet_id` - (Optional) One or more Bare Metal cabinet IDs.
* `epc_host_status` - (Optional) One or more Bare Metal status.
* `host_name` - (Optional) One or more Bare Metal host names.
* `host_type` - (Optional) One or more Bare Metal host types.
* `ids` - (Optional) A list of Bare Metal IDs, all the Bare Metals belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter results by Bare Metal name.
* `os_name` - (Optional) One or more Bare Metal operating system names.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `product_type` - (Optional) One or more Bare Metal product types. valid values: 'lease', 'customer', 'lending'.
* `project_id` - (Optional) One or more project IDs.
* `subnet_id` - (Optional) One or more subnet IDs.
* `vpc_id` - (Optional) One or more vpc IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bare_metals` - It is a nested type which documented below.
  * `availability_zone` - availability zone name.
  * `cabinet_id` - ID of the Cabinet.
  * `cloud_monitor_agent` - The cloud monitor agent of the Bare Metal.
  * `cpu` - cpu specification.
    * `core_count` - number of CPU cores.
    * `count` - number of CPUs.
    * `frequence` - frequence of the cpu.
    * `model` - model of the cpu.
  * `create_time` - The time of creation for Bare Metal.
  * `disk_set` - a list of disks.
    * `disk_type` - type of the disk.
    * `raid` - raid type of the disk.
    * `space` - space of the disk.
  * `enable_bond` - Whether to enable bond.
  * `host_id` - The ID of the Bare Metal.
  * `host_name` - The name of the Bare Metal.
  * `host_status` - status of the Bare Metal.
  * `host_type` - type of the Bare Metal.
  * `image_id` - ID of the Image.
  * `memory` - the memory of the Bare Metal.
  * `network_interface_attribute_set` - a list of network interfaces.
    * `dns1` - DNS1 of the network instance.
    * `dns2` - DNS2 of the network instance.
    * `mac` - MAC of the network instance.
    * `network_interface_id` - the Id of the network interface.
    * `network_interface_type` - type of the network interface.
    * `private_ip_address` - The private IP address assigned to the network interface.
    * `security_group_set` - a list of security groups.
      * `security_group_id` - ID of the security group.
    * `subnet_id` - the ID of the subnet.
    * `vpc_id` - The ID of the VPC.
  * `network_interface_mode` - mode of the network interface.
  * `os_name` - name of the OS.
  * `product_type` - product type of the Bare metal.
  * `raid` - The Raid type of the Bare Metal.
  * `security_agent` - The security agent of the Bare Metal.
  * `sn` - SN of the Bare Metal.
* `total_count` - Total number of Bare Metals that satisfy the condition.


