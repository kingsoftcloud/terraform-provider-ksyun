---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_subnets"
sidebar_current: "docs-ksyun-datasource-subnets"
description: |-
  This data source provides a list of Subnet resources according to their Subnet ID, name and the VPC they belong to.
---

# ksyun_subnets

This data source provides a list of Subnet resources according to their Subnet ID, name and the VPC they belong to.

#

## Example Usage

```hcl
data "ksyun_subnets" "default" {
  output_file            = "output_result"
  ids                    = []
  vpc_id                 = []
  nat_id                 = []
  network_acl_id         = []
  subnet_type            = []
  availability_zone_name = []
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone_names` - (Optional) The availability zone that the desired Subnet belongs to.
* `ids` - (Optional) A list of Subnet IDs, all the Subnet resources belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter results by subnet name.
* `nat_ids` - (Optional) The id of the NAT that the desired Subnet associated to.
* `network_acl_ids` - (Optional) The id of the ACL that the desired Subnet associated to.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `subnet_types` - (Optional) one or more subnet types.
* `vpc_ids` - (Optional) The id of the VPC that the desired Subnet belongs to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `subnets` - An information list of subnets. Each element contains the following attributes:
  * `availability_zone_name` - Availability zone.
  * `availble_i_p_number` - number of available IPs.
  * `cidr_block` - The CIDR block assigned to the subnet.
  * `create_time` - creation time of the subnet.
  * `dhcp_ip_from` - DHCP start IP.
  * `dhcp_ip_to` - DHCP end IP.
  * `dns1` - The dns1 of the subnet.
  * `dns2` - The dns2 of the subnet.
  * `gateway_ip` - The IP of gateway.
  * `id` - ID of the subnet.
  * `ipv6_cidr_block_association_set` - An Ipv6 association list of this vpc.
    * `ipv6_cidr_block` - the Ipv6 of this vpc bound.
  * `name` - Name of the subnet.
  * `nat_id` - The id of the NAT that the desired Subnet associated to.
  * `network_acl_id` - The id of the ACL that the desired Subnet associated to.
  * `subnet_id` - ID of the subnet.
  * `subnet_type` - Type of the subnet.
  * `vpc_id` - ID of the VPC.
* `total_count` - Total number of Subnet resources that satisfy the condition.


