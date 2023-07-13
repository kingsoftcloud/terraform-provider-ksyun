---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_vpc"
sidebar_current: "docs-ksyun-resource-vpc"
description: |-
  Provides a VPC resource.
---

# ksyun_vpc

Provides a VPC resource.

~> **Note**  The network segment can only be created or deleted, can not perform both of them at the same time.

## Example Usage



## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional, ForceNew) The CIDR blocks of VPC.
* `is_default` - (Optional, ForceNew) Whether the VPC is default or not.
* `provided_ipv6_cidr_block` - (Optional, ForceNew) whether support IPV6 CIDR blocks. <br> NOTES: providing a part of regions now.
* `vpc_name` - (Optional) The name of the vpc.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time of creation for VPC.
* `ipv6_cidr_block_association_set` - An Ipv6 association list of this vpc.
  * `ipv6_cidr_block` - the Ipv6 of this vpc bound.


## Import

VPC can be imported using the `id`, e.g.

```
$ terraform import ksyun_vpc.example vpc-abc123456
```

