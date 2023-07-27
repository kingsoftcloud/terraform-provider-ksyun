---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_nat_instance_bandwidth_limit"
sidebar_current: "docs-ksyun-resource-nat_instance_bandwidth_limit"
description: |-
  Provides a bandwidth limit rule of Nat and Instance resource under VPC resource.
---

# ksyun_nat_instance_bandwidth_limit

Provides a bandwidth limit rule of Nat and Instance resource under VPC resource.

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

  #max_count=1
  #min_count=1
  subnet_id         = "${ksyun_subnet.foo.id}"
  instance_password = "Xuan663222"
  keep_image_login  = false
  charge_type       = "Daily"
  purchase_time     = 1
  security_group_id = ["${ksyun_security_group.default.id}"]
  instance_name     = "ksyun-kec-tf-nat"
  sriov_net_support = "false"
}

resource "ksyun_nat" "foo" {
  nat_name    = "ksyun-nat-tf"
  nat_mode    = "Subnet"
  nat_type    = "public"
  band_width  = 10
  charge_type = "DailyPaidByTransfer"
  vpc_id      = "${ksyun_vpc.foo.id}"
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

resource "ksyun_nat_associate" "foo" {
  nat_id               = "${ksyun_nat.foo.id}"
  network_interface_id = "${ksyun_instance.foo.network_interface_id}"
}

resource "ksyun_nat_instance_bandwidth_limit" "foo" {
  nat_id               = "${ksyun_nat.foo.id}"
  network_interface_id = "${ksyun_instance.foo.network_interface_id}"
  bandwidth_limit      = 1
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth_limit` - (Required) The bandwidth limit of network interface that belong to instance.
* `nat_id` - (Required, ForceNew) The id of the Nat.
* `network_interface_id` - (Required, ForceNew) The id of network interface that belong to instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_type` - the type of instance. Values: epc or kec.
* `private_ip_address` - the private ip of network interface.


