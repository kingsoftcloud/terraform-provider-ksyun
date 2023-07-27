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
}

resource "ksyun_nat_associate" "foo" {
  nat_id    = "${ksyun_nat.foo.id}"
  subnet_id = "${ksyun_subnet.test.id}"
}
resource "ksyun_nat_associate" "associate_ins" {
  nat_id               = "${ksyun_nat.foo.id}"
  network_interface_id = "${ksyun_instance.foo.network_interface_id}"
}
```

## Argument Reference

The following arguments are supported:

* `nat_id` - (Required, ForceNew) The id of the Nat.
* `network_interface_id` - (Optional, ForceNew) The id of network interface that belong to instance. Notes: Because of there is one resource in the association, conflict with `subnet_id`.
* `subnet_id` - (Optional, ForceNew) The id of the Subnet. Notes: Because of there is one resource in the association, conflict with `network_interface_id`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

nat associate can be imported using the `id`, the id format must be `{nat_id}:{resource_id}`, resource_id range `subnet_id`, `natwork_interface_id` e.g.

## Import Subnet association
```
$ terraform import ksyun_nat_associate.example $nat_id:subnet-$subnet_id
```
## Import NetworkInterface association
```
$ terraform import ksyun_nat_associate.example $nat_id:kni-$network_interface_id
```

