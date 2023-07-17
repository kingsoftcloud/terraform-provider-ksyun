---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_kec_network_interface"
sidebar_current: "docs-ksyun-resource-kec_network_interface"
description: |-
  Provides a KEC network interface resource.
---

# ksyun_kec_network_interface

Provides a KEC network interface resource.

#

## Example Usage

```hcl
resource "ksyun_kec_network_interface" "default" {
  subnet_id              = "81530211-2785-47a8-b2a0-ae13120fa97d"
  security_group_ids     = ["7e2f45b5-e79d-4612-a7fc-fe74a50b639a", "35ac2642-1958-4ed7-b02c-dc86f27bc9d9"]
  network_interface_name = "Ksc_NetworkInterface"
}
```

## Argument Reference

The following arguments are supported:

* `security_group_ids` - (Required) A list of security group IDs.
* `subnet_id` - (Required) The ID of the subnet which the network interface belongs to.
* `network_interface_name` - (Optional) The name of the network interface.
* `private_ip_address` - (Optional) Private IP.
* `secondary_private_ip_address_count` - (Optional) The count of secondary private id address automatically assigned. <br> Notes:  `secondary_private_ip_address_count` conflict with `secondary_private_ips`.
* `secondary_private_ips` - (Optional) Assign secondary private ips to the network interface. <br> Notes: `secondary_private_ips` conflict with `secondary_private_ip_address_count`.

The `secondary_private_ips` object supports the following:

* `ip` - (Required) Secondary Private IP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_id` - The instance id to bind with the network interface.


## Import

Instance can be imported using the id, e.g.

```
$ terraform import ksyun_kec_network_interface.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

