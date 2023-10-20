---
subcategory: "VPN"
layout: "ksyun"
page_title: "ksyun: ksyun_vpn_customer_gateway"
sidebar_current: "docs-ksyun-resource-vpn_customer_gateway"
description: |-
  Provides a Vpn Customer Gateway resource.
---

# ksyun_vpn_customer_gateway

Provides a Vpn Customer Gateway resource.

#

## Example Usage

```hcl
resource "ksyun_vpn_customer_gateway" "default" {
  customer_gateway_address    = "100.0.0.2"
  ha_customer_gateway_address = "100.0.2.2"
  customer_gateway_name       = "ksyun_vpn_cus_gw"
}
```

## Argument Reference

The following arguments are supported:

* `customer_gateway_address` - (Required) The customer gateway address of the vpn customer gateway.
* `customer_gateway_name` - (Required) The name of the vpn customer gateway.
* `ha_customer_gateway_address` - (Required) The ha customer gateway address of the vpn customer gateway.
* `customer_gateway_mame` - (Optional) The name of the vpn customer gateway. **Warning this field was removed**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Vpn Customer Gateway can be imported using the `id`, e.g.

```
$ terraform import ksyun_vpn_customer_gateway.default $id
```

