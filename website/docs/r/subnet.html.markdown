---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_subnet"
sidebar_current: "docs-ksyun-resource-subnet"
description: |-
  Provides a Subnet resource under VPC resource.
---

# ksyun_subnet

Provides a Subnet resource under VPC resource.

#

## Example Usage

```hcl
resource "ksyun_vpc" "example" {
  vpc_name   = "tf-example-vpc-01"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "example" {
  subnet_name       = "tf-acc-subnet1"
  cidr_block        = "10.0.5.0/24"
  subnet_type       = "Normal"
  dhcp_ip_from      = "10.0.5.2"
  dhcp_ip_to        = "10.0.5.253"
  vpc_id            = "${ksyun_vpc.test.id}"
  gateway_ip        = "10.0.5.1"
  dns1              = "198.18.254.41"
  dns2              = "198.18.254.40"
  availability_zone = "cn-shanghai-2a"
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required, ForceNew) The CIDR block assigned to the subnet.
* `subnet_type` - (Required, ForceNew) The type of subnet. Valid Values:'Reserve', 'Normal', 'Physical'.
* `vpc_id` - (Required, ForceNew) The id of the vpc.
* `availability_zone` - (Optional, ForceNew) The name of the availability zone.
* `dhcp_ip_from` - (Optional, ForceNew, **Deprecated**) This attribute is deprecated and will be removed in a future version. DHCP start IP.
* `dhcp_ip_to` - (Optional, ForceNew, **Deprecated**) This attribute is deprecated and will be removed in a future version. DHCP end IP.
* `dns1` - (Optional) The dns of the subnet.
* `dns2` - (Optional) The dns of the subnet.
* `gateway_ip` - (Optional, ForceNew) The IP of gateway.
* `provided_ipv6_cidr_block` - (Optional, ForceNew) whether support IPV6 CIDR blocks. <br> NOTES: providing a part of regions now.
* `subnet_name` - (Optional) The name of the subnet.
* `visit_internet` - (Optional) Whether the subnet can access the Internet. Valid, when subnet_type = Physical.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `availability_zone_name` - The name of the availability zone.
* `available_ip_number` - number of available IPs.
* `create_time` - creation time of the subnet.
* `ipv6_cidr_block_association_set` - An Ipv6 association list of this subnet.
  * `ipv6_cidr_block` - the Ipv6 of this subnet bound.
* `nat_id` - The id of the NAT that the desired Subnet associated to.
* `network_acl_id` - The id of the ACL that the desired Subnet associated to.
* `subnet_id` - ID of the subnet.


## Import

Subnet can be imported using the `id`, e.g.

```
$ terraform import ksyun_subnet.example fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```

