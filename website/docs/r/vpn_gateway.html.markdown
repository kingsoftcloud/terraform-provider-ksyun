---
subcategory: "VPN"
layout: "ksyun"
page_title: "ksyun: ksyun_vpn_gateway"
sidebar_current: "docs-ksyun-resource-vpn_gateway"
description: |-
  Provides a Vpn Gateway resource under VPC resource.
---

# ksyun_vpn_gateway

Provides a Vpn Gateway resource under VPC resource.

#

## Example Usage

```hcl
# create vpn gateway with vpn 1.0 version
resource "ksyun_vpn_gateway" "default" {
  vpn_gateway_name = "ksyun_vpn_gw_tf1"
  band_width       = 10
  vpc_id           = "a8979fe2-cf1a-47b9-80f6-57445227c541"
  charge_type      = "Daily"
  # vpn_gateway_version = "1.0"
}

# create vpn gateway with vpn 2.0 version
resource "ksyun_vpn_gateway" "default" {
  vpn_gateway_name    = "ksyun_vpn_gw_tf1"
  band_width          = 10
  vpc_id              = "a8979fe2-cf1a-47b9-80f6-57445227c541"
  charge_type         = "Daily"
  vpn_gateway_version = "2.0"
}
```

## Argument Reference

The following arguments are supported:

* `band_width` - (Required) The bandWidth of the vpn gateway.Valid Values:5,10,20,50,100,200.
* `charge_type` - (Required, ForceNew) The charge type of the vpn gateway.Valid Values:'Monthly','Daily'.
* `vpc_id` - (Required, ForceNew) The id of the vpc.
* `project_id` - (Optional) The project id  of the vpn gateway.Default is 0.
* `purchase_time` - (Optional, ForceNew) The purchase time of the vpn gateway.
* `vpn_gateway_name` - (Optional) The name of the vpn gateway.
* `vpn_gateway_version` - (Optional, ForceNew) the version of vpn gateway. Default `1.0`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `vpn_gateway_id` - The ID of the vpn gateway.


## Import

Vpn Gateway can be imported using the `id`, e.g.

```
$ terraform import ksyun_vpn_gateway.default $id
```

