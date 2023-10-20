---
subcategory: "VPN"
layout: "ksyun"
page_title: "ksyun: ksyun_vpn_gateway_route"
sidebar_current: "docs-ksyun-resource-vpn_gateway_route"
description: |-
  Provides a Vpn Gateway Route resource under VPC resource.
**Notes:** `ksyun_vpn_gateway_route` only valid when Vpn 2.0
---

# ksyun_vpn_gateway_route

Provides a Vpn Gateway Route resource under VPC resource.
**Notes:** `ksyun_vpn_gateway_route` only valid when Vpn 2.0

#

## Example Usage

```hcl
resource "ksyun_vpn_gateway_route" "default1" {
  vpn_gateway_id         = "450a71b0-ea20-****-*****"
  next_hop_type          = "vpc"
  destination_cidr_block = "10.7.255.0/30"
}
```

## Argument Reference

The following arguments are supported:

* `destination_cidr_block` - (Required, ForceNew) The destination cidr block.
* `next_hop_type` - (Required, ForceNew) The type of next hop. Valid Values: `vpn_tunnel`, `vpc`.
* `vpn_gateway_id` - (Required, ForceNew) The ID of the vpn gateway.
* `next_hop_instance_id` - (Optional, ForceNew) The instance id of next hop, which must be set when `next_hop_type` is `vpn_tunnel.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



