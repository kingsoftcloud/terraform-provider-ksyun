---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_private_dns_zone_vpc_attachment"
sidebar_current: "docs-ksyun-resource-private_dns_zone_vpc_attachment"
description: |-
  Provides a resource to create a PrivateDns zone_vpc_attachment
---

# ksyun_private_dns_zone_vpc_attachment

Provides a resource to create a PrivateDns zone_vpc_attachment

#

## Example Usage

```hcl
provider "ksyun" {
  region = "cn-guangzhou-1"
}

resource "ksyun_private_dns_zone" "foo" {
  zone_name   = "tf-pdns-binding.com"
  zone_ttl    = 360
  charge_type = "TrafficMonthly"
}

resource "ksyun_vpc" "foo" {
  vpc_name   = "tf-pdns-binding-vpc"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_private_dns_zone_vpc_attachment" "example" {
  zone_id = ksyun_private_dns_zone.foo.id
  vpc_set {
    region_name = "cn-guangzhou-1"
    vpc_id      = ksyun_vpc.foo.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Private Dns Zone ID.
* `vpc_set` - (Optional, ForceNew) New add vpc info.

The `vpc_set` object supports the following:

* `region_name` - (Required) Vpc region.
* `vpc_id` - (Required) Vpc Id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Private Dns zone_vpc_attachment can be imported using the id, e.g.

```
terraform import ksyun_private_dns_zone_vpc_attachment.example ${zone_id}:${vpc_id}
```

