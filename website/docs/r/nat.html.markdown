---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_nat"
sidebar_current: "docs-ksyun-resource-nat"
description: |-
  Provides a Nat resource under VPC resource.
---

# ksyun_nat

Provides a Nat resource under VPC resource.

#

## Example Usage

```hcl
resource "ksyun_vpc" "test" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_nat" "foo" {
  nat_name    = "ksyun-nat-tf"
  nat_mode    = "Vpc"
  nat_type    = "public"
  band_width  = 1
  charge_type = "DailyPaidByTransfer"
  vpc_id      = "${ksyun_vpc.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `band_width` - (Required) The BandWidth of Nat Ip, value range:[1, 15000], Default is 1.
* `nat_mode` - (Required, ForceNew) Mode of the NAT, valid values: 'Vpc', 'Subnet'.
* `vpc_id` - (Required, ForceNew) ID of the VPC.
* `charge_type` - (Optional, ForceNew) charge type, valid values: 'Monthly', 'Peak', 'Daily', 'PostPaidByAdvanced95Peak', 'DailyPaidByTransfer'. Default is DailyPaidByTransfer.
* `nat_ip_number` - (Optional) The Counts of Nat Ip, value range:[1, 20], Default is 1.
* `nat_line_id` - (Optional) ID of the line.
* `nat_name` - (Optional) Name of the NAT.
* `nat_type` - (Optional, ForceNew) Type of the NAT, valid values: 'public'.
* `project_id` - (Optional) ID of the project.
* `purchase_time` - (Optional, ForceNew) The PurchaseTime of the Nat, value range [1, 36]. If charge_type is Monthly this Field is Required.
* `tags` - (Optional) the tags of the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time of creation of Nat.
* `nat_ip_set` - The nat ip list of the desired Nat.
  * `nat_ip_id` - The ID of the NAT IP.
  * `nat_ip` - NAT IP address.


## Import

nat can be imported using the `id`, e.g.

```
$ terraform import ksyun_nat.example fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```

