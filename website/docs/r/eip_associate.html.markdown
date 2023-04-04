---
subcategory: "EIP"
layout: "ksyun"
page_title: "ksyun: ksyun_eip_associate"
sidebar_current: "docs-ksyun-resource-eip_associate"
description: |-
  Provides an EIP Association resource for associating Elastic IP to UHost Instance, Load Balancer, etc.
---

# ksyun_eip_associate

Provides an EIP Association resource for associating Elastic IP to UHost Instance, Load Balancer, etc.

## Example Usage

```hcl
resource "ksyun_eip_associate" "slb" {
  allocation_id = "419782b7-6766-4743-afb7-7c7081214092"
  instance_type = "Slb"
  instance_id   = "7fae85e4-ab1a-415c-aef9-03a402c79d97"
}
resource "ksyun_eip_associate" "server" {
  allocation_id        = "419782b7-6766-4743-afb7-7c7081214092"
  instance_type        = "Ipfwd"
  instance_id          = "566567677-6766-4743-afb7-7c7081214092"
  network_interface_id = "87945980-59659-04548-759045803"
}
```

## Argument Reference

The following arguments are supported:

* `allocation_id` - (Required, ForceNew) The ID of EIP.
* `instance_id` - (Required, ForceNew) The id of the instance.
* `instance_type` - (Required, ForceNew) The type of the instance.Valid Values:'Ipfwd', 'Slb'.
* `network_interface_id` - (Optional, ForceNew) The id of the network interface.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `band_width_share_id` - the ID of the BWS which the EIP associated.
* `band_width` - The band width of the public address.
* `create_time` - creation time of the EIP.
* `internet_gateway_id` - InternetGateway ID.
* `ip_version` - IP version of the EIP.
* `is_band_width_share` - BWS EIP.
* `line_id` - The id of the line.
* `project_id` - The id of the project.
* `public_ip` - The Elastic IP address.
* `state` - state of the EIP.


## Import

EIP Association can be imported using the id, e.g.

```
$ terraform import ksyun_eip_associate.default ${allocation_id}:${instance_id}:${network_interface_id}
```

