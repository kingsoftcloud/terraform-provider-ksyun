---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_route"
sidebar_current: "docs-ksyun-resource-route"
description: |-
  Provides a route resource under VPC resource.
---

# ksyun_route

Provides a route resource under VPC resource.

#

## Example Usage

```hcl
resource "ksyun_vpc" "example" {
  vpc_name   = "tf-example-vpc-01"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_route" "example" {
  destination_cidr_block = "10.0.0.0/16"
  route_type             = "InternetGateway"
  vpc_id                 = "${ksyun_vpc.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `destination_cidr_block` - (Required, ForceNew) The CIDR block assigned to the route.
* `route_type` - (Required, ForceNew) The type of route.Valid Values:'InternetGateway', 'Tunnel', 'Host', 'Peering', 'DirectConnect', 'Vpn'.
* `vpc_id` - (Required, ForceNew) The id of the vpc.
* `direct_connect_gateway_id` - (Optional, ForceNew) The id of the DirectConnectGateway, If route_type is DirectConnect, This Field is Required.
* `instance_id` - (Optional, ForceNew) The id of the VM, If route_type is Host, This Field is Required.
* `tunnel_id` - (Optional, ForceNew) The id of the tunnel If route_type is Tunnel, This Field is Required.
* `vpc_peering_connection_id` - (Optional, ForceNew) The id of the Peering, If route_type is Peering, This Field is Required.
* `vpn_tunnel_id` - (Optional, ForceNew) The id of the Vpn, If route_type is Vpn, This Field is Required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time of creation of the route.
* `next_hop_set` - A list of next hop.
  * `gateway_id` - The ID of the gateway.
  * `gateway_name` - The name of the gateway.


## Import

route can be imported using the `id`, e.g.

```
$ terraform import ksyun_route.example xxxx-xxxxx
```

