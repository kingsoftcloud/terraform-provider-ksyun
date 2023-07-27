---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_dnat"
sidebar_current: "docs-ksyun-resource-dnat"
description: |-
  Provides a Dnat resource under VPC resource.
---

# ksyun_dnat

Provides a Dnat resource under VPC resource.

#

## Example Usage

```hcl
data "ksyun_images" "centos-7_5" {
  platform = "centos-7.5"
}
data "ksyun_availability_zones" "default" {
}

resource "ksyun_security_group" "default" {
  vpc_id              = "${ksyun_vpc.foo.id}"
  security_group_name = "ksyun-security-group-nat"
}

resource "ksyun_instance" "foo" {
  image_id      = "${data.ksyun_images.centos-7_5.images.0.image_id}"
  instance_type = "N3.2B"

  subnet_id         = "${ksyun_subnet.foo.id}"
  instance_password = "Xuan663222"
  keep_image_login  = false
  charge_type       = "Daily"
  purchase_time     = 1
  security_group_id = ["${ksyun_security_group.default.id}"]
  instance_name     = "ksyun-kec-tf-nat"
  sriov_net_support = "false"
  project_id        = 100012
}

resource "ksyun_nat" "foo" {
  nat_name      = "ksyun-nat-tf"
  nat_mode      = "Subnet"
  nat_type      = "public"
  band_width    = 1
  charge_type   = "DailyPaidByTransfer"
  vpc_id        = "${ksyun_vpc.foo.id}"
  nat_ip_number = 3
}
resource "ksyun_vpc" "foo" {
  vpc_name   = "tf-vpc-nat"
  cidr_block = "10.0.5.0/24"
}

resource "ksyun_subnet" "foo" {
  subnet_name       = "tf-acc-nat-subnet1"
  cidr_block        = "10.0.5.0/24"
  subnet_type       = "Normal"
  vpc_id            = "${ksyun_vpc.foo.id}"
  gateway_ip        = "10.0.5.1"
  dns1              = "198.18.254.41"
  dns2              = "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}

resource "ksyun_dnat" "foo" {
  nat_id             = ksyun_nat.foo.id
  dnat_name          = "tf_dnat"
  ip_protocol        = "Any"
  nat_ip             = ksyun_nat.foo.nat_ip_set[0].nat_ip
  private_ip_address = ksyun_instance.foo.private_ip_address
  description        = "test"
  private_port       = "Any"
  public_port        = "Any"
}
```

## Argument Reference

The following arguments are supported:

* `ip_protocol` - (Required, ForceNew) the protocol of dnat port, Valid Options: `Any`, `TCP` and `UDP`. <br> Notes: `public_port` and `private_port` must be set as `Any`, when `ip_protocol` is `Any`. Instead, you should set ports.
* `nat_id` - (Required, ForceNew) The id of the Nat.
* `nat_ip` - (Required) the nat ip of nat.
* `private_ip_address` - (Required) the private ip of instance in the identical vpc.
* `description` - (Optional) the description of this dnat rule.
* `dnat_name` - (Optional) the name of dnat rule.
* `private_port` - (Optional) the private port that be accessed in vpc.
* `public_port` - (Optional) the public port of the internet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `dnat_id` - The id of the Subnet. Notes: Because of there is one resource in the association, conflict with `network_interface_id`.


## Import

dnat rules can be imported using the `id`, e.g.

```
$ terraform import ksyun_dnat.foo $dnat_id
```

