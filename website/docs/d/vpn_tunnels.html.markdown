---
subcategory: "VPN"
layout: "ksyun"
page_title: "ksyun: ksyun_vpn_tunnels"
sidebar_current: "docs-ksyun-datasource-vpn_tunnels"
description: |-
  This data source provides a list of VPN tunnels.
---

# ksyun_vpn_tunnels

This data source provides a list of VPN tunnels.

#

## Example Usage

```hcl
data "ksyun_vpn_tunnels" "default" {
  output_file = "output_result"
  ids         = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of VPN tunnel IDs, all the resources belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter results by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `vpn_gateway_ids` - (Optional) A list of vpn gateway ids.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total_count` - Total number of resources that satisfy the condition.
* `vpn_tunnels` - It is a nested type which documented below.
  * `create_time` - creation time.
  * `customer_gateway_id` - Customer gateway ID.
  * `customer_gre_ip` - Customer gre IP.
  * `customer_peer_ip` - The peer ip of customer.
  * `enable_nat_traversal` - The switch of nat traversal.
  * `extra_cidr_set` - A list of extra cidr.
    * `cidr_block` - cidr block.
  * `ha_customer_gre_ip` - HA Customer gre IP.
  * `ha_mode` - The high-availability mode of vpn tunnel.
  * `ha_vpn_gre_ip` - HA VPN gre IP.
  * `health_check_local_peer_cider` - the local peer cider of health checking.
  * `health_check_remote_peer_cider` - The remote peer cider of health checking.
  * `id` - VPN tunnel ID.
  * `ike_authen_algorithm` - IKE authen algorithm.
  * `ike_dh_group` - IKE dh group.
  * `ike_encry_algorithm` - IKE encry algorithm.
  * `ike_version` - The version of ike.
  * `ipsec_authen_algorithm` - IPsec authen algorithm.
  * `ipsec_encry_algorithm` - IPsec encry algorithm.
  * `ipsec_lifetime_second` - IPsec lifetime second.
  * `ipsec_lifetime_traffic` - IPsec lifetime traffic.
  * `local_peer_ip` - The peer ip of kingsoft cloud.
  * `name` - VPN tunnel name.
  * `open_health_check` - The switch of health check.
  * `pre_shared_key` - pre shared key.
  * `state` - VPN tunnel state.
  * `type` - VPN tunnel type.
  * `vpn_gateway_id` - VPN gateway ID.
  * `vpn_gateway_version` - The VPN gateway version.
  * `vpn_gre_ip` - VPN gre IP.
  * `vpn_gwl_name` - The VPN gateway name.
  * `vpn_m_tunnel_create_time` - The creation time of first-vpn-tunnel.
  * `vpn_m_tunnel_state` - The state of master-vpn-tunnel.
  * `vpn_s_tunnel_create_time` - The creation time of second-vpn-tunnel.
  * `vpn_s_tunnel_state` - The state of second-vpn-tunnel.
  * `vpn_tunnel_create_time` - creation time.
  * `vpn_tunnel_id` - VPN tunnel ID.
  * `vpn_tunnel_name` - VPN tunnel name.


