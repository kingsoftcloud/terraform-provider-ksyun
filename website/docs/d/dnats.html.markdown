---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_dnats"
sidebar_current: "docs-ksyun-datasource-dnats"
description: |-
  Query ksyun dnats information
---

# ksyun_dnats

Query ksyun dnats information

#

## Example Usage

```hcl
data "ksyun_dnats" "default" {
  private_ip_address = "10.7.x.xxx"
  nat_id             = "5c7b7925-xxxx-xxxx-xxxx-434fc8042329"
  dnat_ids           = ["5cxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx"]
  output_file        = "output_result"
}
```

## Argument Reference

The following arguments are supported:

* `dnat_ids` - (Optional) The id list of dnats.
* `dnat_name` - (Optional) The name of dnat.
* `ip_protocol` - (Optional) The protocol of dnat rule.
* `nat_id` - (Optional) The nat id of dnat associated.
* `nat_ip` - (Optional) The nat ip.
* `network_interface_id` - (Optional) The network interface id of dnat rule associated.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `private_ip_address` - (Optional) The private ip address.
* `public_port` - (Optional) The public port.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dnats` - An information list of krds db parameter groups. Each element contains the following attributes:
  * `create_time` - The time created.
  * `description` - The description of dnat.
  * `dnat_id` - The id of dnat.
  * `dnat_name` - The name of dnat.
  * `ip_protocol` - The ip protocol of dnat.
  * `nat_id` - The nat id.
  * `nat_ip` - The nat ip of nat associated.
  * `private_ip_address` - The private ip address.
  * `private_port` - The private port.
  * `public_port` - The public port.
* `total_count` - Total number of snapshot policies resources that satisfy the condition.


