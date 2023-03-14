---
subcategory: "VPN"
layout: "ksyun"
page_title: "ksyun: ksyun_vpn_tunnel"
sidebar_current: "docs-ksyun-resource-vpn_tunnel"
description: |-
  Provides a Vpn Tunnel resource.
---

# ksyun_vpn_tunnel

Provides a Vpn Tunnel resource.

#

## Example Usage

```hcl
resource "ksyun_vpn_tunnel" "default" {
  vpn_tunnel_name     = "ksyun_vpn_tunnel_tf_1"
  type                = "Ipsec"
  vpn_gateway_id      = "9b3d361e-f65b-464b-947a-fafb5cfb10d2"
  customer_gateway_id = "7f5a5c91-4814-41bf-b9d6-d9d811f4df0f"
  ike_dh_group        = 2
  pre_shared_key      = "123456789abcd"
}
```

## Argument Reference

The following arguments are supported:

* `customer_gateway_id` - (Required, ForceNew) The customer_gateway_id of the vpn tunnel.
* `pre_shared_key` - (Required, ForceNew) The pre_shared_key of the vpn tunnel.
* `type` - (Required, ForceNew) The bandWidth of the vpn tunnel.Valid Values:'GreOverIpsec','Ipsec'.
* `vpn_gateway_id` - (Required, ForceNew) The vpn_gateway_id of the vpn tunnel.
* `customer_gre_ip` - (Optional, ForceNew) The customer_gre_ip of the vpn tunnel.If type is GreOverIpsec,Required.
* `ha_customer_gre_ip` - (Optional, ForceNew) The ha_customer_gre_ip of the vpn tunnel.If type is GreOverIpsec,Required.
* `ha_vpn_gre_ip` - (Optional, ForceNew) The ha_vpn_gre_ip of the vpn tunnel.If type is GreOverIpsec,Required.
* `ike_authen_algorithm` - (Optional, ForceNew) The ike_authen_algorithm of the vpn tunnel.Valid Values:'md5','sha'.
* `ike_dh_group` - (Optional, ForceNew) The ike_dh_group of the vpn tunnel.Valid Values:1,2,5.
* `ike_encry_algorithm` - (Optional, ForceNew) The ike_encry_algorithm of the vpn tunnel.Valid Values:'3des','aes','des'.
* `ipsec_authen_algorithm` - (Optional, ForceNew) The ipsec_authen_algorithm of the vpn tunnel.Valid Values:'esp-md5-hmac','esp-sha-hmac'.
* `ipsec_encry_algorithm` - (Optional, ForceNew) The ipsec_encry_algorithm of the vpn tunnel.Valid Values:'esp-3des','esp-aes','esp-des','esp-null','esp-seal'.
* `ipsec_lifetime_second` - (Optional, ForceNew) The ipsec_lifetime_second of the vpn tunnel.
* `ipsec_lifetime_traffic` - (Optional, ForceNew) The ipsec_lifetime_traffic of the vpn tunnel.
* `vpn_gre_ip` - (Optional, ForceNew) The vpn_gre_ip of the vpn tunnel.If type is GreOverIpsec, Required.
* `vpn_tunnel_name` - (Optional) The name of the vpn tunnel.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Vpn Tunnel can be imported using the `id`, e.g.

```
$ terraform import ksyun_vpn_tunnel.default $id
```

