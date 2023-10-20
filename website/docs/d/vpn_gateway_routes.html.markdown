---
subcategory: "VPN"
layout: "ksyun"
page_title: "ksyun: ksyun_vpn_gateway_routes"
sidebar_current: "docs-ksyun-datasource-vpn_gateway_routes"
description: |-
  This data source provides a list of VPN GatewayRoutes.
---

# ksyun_vpn_gateway_routes

This data source provides a list of VPN GatewayRoutes.

#

## Example Usage

```hcl
data "ksyun_vpn_gateway_routes" "default" {
  output_file = "output_result"

  # specify vpn_gateway_id to query vpn_gateway_routes
  vpn_gateway_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `cidr_blocks` - (Optional) A list of cidr block.
* `next_hop_types` - (Optional) A list of the next hop type.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `vpn_gateway_id` - (Optional) A list of VPN gateway IDs, all the resources belong to this region will be retrieved if the ID is `""`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total_count` - Total number of resources that satisfy the condition.
* `vpn_gateway_routes` - It is a nested type which documented below.
  * `create_time` - Creation time.
  * `destination_cidr_block` - The ID of the gateway.
  * `next_hop_instance_name` - Band width.
  * `next_hop_type` - The name of the gateway.
  * `route_type` - The name of the gateway.
  * `vpn_gateway_id` - VPC ID.
  * `vpn_gateway_route_id` - The ID of the gateway.


