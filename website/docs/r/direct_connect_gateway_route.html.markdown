---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_direct_connect_gateway_route"
sidebar_current: "docs-ksyun-resource-direct_connect_gateway_route"
description: |-
  Provides a DirectConnectGatewayRoute resource.
---

# ksyun_direct_connect_gateway_route

Provides a DirectConnectGatewayRoute resource.

## Example Usage

```hcl
data "ksyun_direct_connects" "test" {
  ids        = []
  name_regex = ".*test.*"
}

resource "ksyun_direct_connect_interface" "test" {
  direct_connect_id  = data.ksyun_direct_connects.test.direct_connects[0].id
  route_type         = "STATIC"
  bgp_peer           = 59019
  bgp_client_token   = "dadasd"
  reliability_method = "bfd"
  enable_ipv6        = true
  bfd_config_id      = "29e0c675-2cca-4778-b331-884fca06de17"
  vlan_id            = 111

  direct_connect_interface_name = "tf_direct_connect_test_1"
}

resource "ksyun_direct_connect_gateway" "test" {
  direct_connect_gateway_name = "tf_direct_connect_gateway_test_1"
  vpc_id                      = "a38673ae-c9b7-4f8e-b727-b6feb648805b"
}

resource "ksyun_direct_connect_gateway_route" "test" {
  direct_connect_gateway_id = ksyun_direct_connect_gateway.test.id
  destination_cidr_block    = "192.136.0.0/24"
  next_hop_type             = "Vpc"
  depends_on                = [ksyun_dc_interface_associate.test]
}
```

## Argument Reference

The following arguments are supported:

* `destination_cidr_block` - (Required, ForceNew) The destination CIDR block of the route. The CIDR block must be in the format of `x.x.x.x/x`.
* `direct_connect_gateway_id` - (Required, ForceNew) The ID of the direct connect gateway.
* `next_hop_type` - (Required, ForceNew) The type of the next hop. Valid values: `Vpc`, `DirectConnect`, `Cen`. Default is `Vpc`. If set to `DirectConnect`, the next hop instance ID must be provided.
* `bgp_status` - (Optional) BGP Status.
* `enable_ip_v6` - (Optional) whether to enable IPv6. Valid values: `true`, `false`. Default is `false`.
* `next_hop_instance` - (Optional) The next hop instance ID.
* `priority` - (Optional) Priority.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `as_path` - AS Path of the route.
* `direct_connect_gateway_route_id` - The ID of the direct connect interface.
* `direct_connect_id` - Direct Connect ID.
* `next_hop_instance_name` - The name of the next hop instance.
* `route_type` - Route Type.


## Import

Route can be imported using the id, e.g.

```
$ terraform import ksyun_direct_connect_gateway_route.test 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

