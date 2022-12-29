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
* `vpc_name` - (Optional) The name of the vpc.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time of creation for VPC, formatted in RFC3339 time string.


## Import

VPC can be imported using the `id`, e.g.

```
$ terraform import ksyun_vpc.example vpc-abc123456
```

