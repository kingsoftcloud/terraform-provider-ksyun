---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_network_interfaces"
sidebar_current: "docs-ksyun-datasource-network_interfaces"
description: |-
  This data source provides a list of Network Interface resources according to their Network Interface ID.
---

# ksyun_network_interfaces

This data source provides a list of Network Interface resources according to their Network Interface ID.

#

## Example Usage

```hcl
data "ksyun_network_interfaces" "default" {
  output_file        = "output_result"
  ids                = []
  securitygroup_id   = []
  instance_type      = []
  instance_id        = []
  private_ip_address = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Network Interface IDs, all the Network Interfaces belong to this region will be retrieved if the ID is `""`.
* `instance_id` - (Optional) A list of VPC instance IDs.
* `instance_type` - (Optional) A list of instance types.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `private_ip_address` - (Optional) A list of private IPs.
* `securitygroup_id` - (Optional) A list of security group IDs.
* `subnet_id` - (Optional) A list of subnet IDs.
* `vpc_id` - (Optional) A list of VPC IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `network_interfaces` - An information list of network interfaces. Each element contains the following attributes:
  * `assigned_private_ip_address_set` - Assign secondary private ips to the network interface.
  * `d_n_s1` - DNS 1.
  * `d_n_s2` - DNS 2.
  * `id` - The id of the network interface.
  * `instance_id` - The ID of the instance.
  * `instance_type` - The type of the instance.
  * `mac_address` - The mac address of the network interface.
  * `name` - The name of the network interface.
  * `network_interface_id` - The id of the network interface.
  * `network_interface_name` - The name of the network interface.
  * `network_interface_type` - The type of the network interface.
  * `private_ip_address` - private IP.
  * `security_group_set` - A list of security groups.
    * `security_group_id` - The ID of the security group.
    * `security_group_name` - The name of the security group.
* `total_count` - Total number of network interface resources that satisfy the condition.


