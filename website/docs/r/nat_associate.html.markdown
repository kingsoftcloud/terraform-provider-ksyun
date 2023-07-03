---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_nat_associate"
sidebar_current: "docs-ksyun-resource-nat_associate"
description: |-
  Provides a Nat Associate resource under VPC resource.
---

# ksyun_nat_associate

Provides a Nat Associate resource under VPC resource.

#

## Example Usage

```hcl
resource "ksyun_vpc" "test" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_nat" "foo" {
  nat_name    = "ksyun-nat-tf"
  nat_mode    = "Subnet"
  nat_type    = "public"
  band_width  = 1
  charge_type = "DailyPaidByTransfer"
  vpc_id      = "${ksyun_vpc.test.id}"
}

resource "ksyun_subnet" "test" {
  subnet_name       = "tf-acc-subnet1"
  cidr_block        = "10.0.5.0/24"
  subnet_type       = "Normal"
  dhcp_ip_from      = "10.0.5.2"
  dhcp_ip_to        = "10.0.5.253"
  vpc_id            = "${ksyun_vpc.test.id}"
  gateway_ip        = "10.0.5.1"
  dns1              = "198.18.254.41"
  dns2              = "198.18.254.40"
  availability_zone = "cn-beijing-6a"
}

resource "ksyun_nat_associate" "foo" {
  nat_id    = "${ksyun_nat.foo.id}"
  subnet_id = "${ksyun_subnet.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `nat_id` - (Required, ForceNew) The id of the Nat.
* `subnet_id` - (Required, ForceNew) The id of the Subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

nat associate can be imported using the `id`, e.g.

```
$ terraform import ksyun_nat_associate.example $nat_id:$subnet_id
```

