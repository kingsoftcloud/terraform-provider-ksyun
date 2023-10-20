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
# create Vpn Tunnel with Vpn 1.0
resource "ksyun_vpn_tunnel" "tunnel-vpn1" {
  vpn_tunnel_name     = "tf_vpn_tunnel_vpn1"
  type                = "Ipsec"
  vpn_gateway_id      = "9b3d361e-f65b-464b-947a-fafb5cfb10d2"
  customer_gateway_id = "7f5a5c91-4814-41bf-b9d6-d9d811f4df0f"
  ike_dh_group        = 2
  pre_shared_key      = "123456789abcd"
}

# create Vpn Tunnel with Vpn 2.0
resource "ksyun_vpn_tunnel" "tunnel-vpn2" {
  vpn_gateway_version = "2.0" # choose vpn gateway version
  vpn_tunnel_name     = "tf_vpn_tunnel_vpn2"
  type                = "Ipsec"
  ike_version         = "v1"
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
* `type` - (Required, ForceNew) The bandWidth of the vpn tunnel. Valid Values: VPN-v1: 'GreOverIpsec' or 'Ipsec'; VPN-v2: `RouteIpsec` or `Ipsec`.
* `vpn_gateway_id` - (Required, ForceNew) The vpn_gateway_id of the vpn tunnel.
* `customer_gre_ip` - (Optional, ForceNew) The customer_gre_ip of the vpn tunnel.If type is GreOverIpsec and Vpn-Gateway-Version is 1.0, Required. Notes: it's valid when vpn gateway version is 1.0.
* `customer_peer_ip` - (Optional) The IP of customer with CIDR indicated. Notes: it's valid when vpn gateway version is 2.0.
* `ha_customer_gre_ip` - (Optional, ForceNew) The ha_customer_gre_ip of the vpn tunnel.If type is GreOverIpsec and Vpn-Gateway-Version is 1.0, Required. Notes: it's valid when vpn gateway version is 1.0.
* `ha_mode` - (Optional, ForceNew) The high-availability mode of vpn tunnel. Valid values: `active_active` valid only when type as `Ipsec`; `active_active` and `active_standby` valid only when type as `RouteIpsec`. Notes: it's valid when vpn gateway version is 2.0.
* `ha_vpn_gre_ip` - (Optional, ForceNew) The ha_vpn_gre_ip of the vpn tunnel.If type is GreOverIpsec,Required and Vpn-Gateway-Version is 1.0, Required. Notes: it's valid when vpn gateway version is 1.0.
* `ike_authen_algorithm` - (Optional, ForceNew) The ike_authen_algorithm of the vpn tunnel.Valid Values:'md5','sha'.
* `ike_dh_group` - (Optional, ForceNew) The ike_dh_group of the vpn tunnel.Valid Values:1,2,5.
* `ike_encry_algorithm` - (Optional, ForceNew) The ike_encry_algorithm of the vpn tunnel.Valid Values:'3des','aes','des'.
* `ike_version` - (Optional, ForceNew) the version of Ike. Notes: it's valid when vpn gateway version is 2.0.
* `ipsec_authen_algorithm` - (Optional, ForceNew) The ipsec_authen_algorithm of the vpn tunnel.Valid Values:'esp-md5-hmac','esp-sha-hmac'.
* `ipsec_encry_algorithm` - (Optional, ForceNew) The ipsec_encry_algorithm of the vpn tunnel.Valid Values:'esp-3des','esp-aes','esp-des','esp-null','esp-seal'.
* `ipsec_lifetime_second` - (Optional, ForceNew) The ipsec_lifetime_second of the vpn tunnel.
* `ipsec_lifetime_traffic` - (Optional, ForceNew) The ipsec_lifetime_traffic of the vpn tunnel.
* `local_peer_ip` - (Optional) The local IP in Kingsoft Cloud with CIDR indicated. Notes: it's valid when vpn gateway version is 2.0.
* `open_health_check` - (Optional, ForceNew) The switch of vpn tunnel health check. **Notes: that's valid only when vpn-v2.0 and tunnel type is `RouteIpsec`**. Notes: it's valid when vpn gateway version is 2.0.
* `vpn_gateway_version` - (Optional, ForceNew) The version of vpn gateway. The version must be identical with `vpn_gate_way_version` of `ksyun_vpn_gateway`.
* `vpn_gre_ip` - (Optional, ForceNew) The vpn_gre_ip of the vpn tunnel. If type is GreOverIpsec and Vpn-Gateway-Version is 1.0, Required. Notes: it's valid when vpn gateway version is 1.0.
* `vpn_tunnel_name` - (Optional) The name of the vpn tunnel.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `vpn_m_tunnel_create_time` - the vpn first tunnel created time.
* `vpn_m_tunnel_state` - the vpn first tunnel state.
* `vpn_s_tunnel_create_time` - the vpn second tunnel created time.
* `vpn_s_tunnel_state` - the vpn second tunnel state.
* `vpn_tunnel_create_time` - the vpn tunnel created time.


## Import

Vpn Tunnel can be imported using the `id`, e.g.

```
$ terraform import ksyun_vpn_tunnel.default $id
```

