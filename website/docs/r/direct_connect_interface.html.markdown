---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_direct_connect_interface"
sidebar_current: "docs-ksyun-resource-direct_connect_interface"
description: |-
  Provides a DirectConnectInterface resource.
---

# ksyun_direct_connect_interface

Provides a DirectConnectInterface resource.

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
```

## Argument Reference

The following arguments are supported:

* `direct_connect_id` - (Required, ForceNew) The id of direct connect. It's meaning is the physical port.
* `bfd_config_id` - (Optional) The ID of the BFD configuration.
* `bgp_client_token` - (Optional) Bgp client token is used to ensure the idempotency of the request. It can be any string, but it must be unique for each request.
* `bgp_peer` - (Optional) The BGP peer IP address. It is used to establish a BGP session with the customer.
* `customer_ipv6_peer_ip` - (Optional) Customer IPv6 peer IP address.
* `customer_peer_ip` - (Optional) Customer peer IP address. It is used to establish a BGP session with the customer.
* `direct_connect_interface_account_id` - (Optional) The account ID of the direct connect interface. It is used to create a direct connect interface in another account.
* `direct_connect_interface_name` - (Optional) The name of the direct connect interface. It is used to identify the direct connect interface.
* `enable_ipv6` - (Optional) Enable IPv6. Valid values: `true`, `false`. Default is `false`.
* `ha_customer_peer_ip` - (Optional) Ha customer peer IP address.
* `ha_direct_connect_id` - (Optional) Ha direct connect ID. It is used to create a high availability direct connect interface.
* `ha_local_peer_ip` - (Optional) Ha customer peer IP address.
* `local_ipv6_peer_ip` - (Optional) Local IPv6 peer IP address.
* `local_peer_ip` - (Optional) Local peer IP address. It is used to establish a BGP session with the customer.
* `reliability_method` - (Optional) Reliability method. Valid values: `bfd`, `nqa`. Default is `nqa`. If set to `BFD`, BFD configuration must be provided.
* `route_type` - (Optional, ForceNew) Route Type. Valid values: `BGP`, `STATIC`. Default is `BGP`. If set to `STATIC`, the customer must provide the BGP peer IP address and local peer IP address.
* `vlan_id` - (Optional) The id of vlan in direct connect.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `account_id` - Account ID of the direct connect interface.
* `customer_peer_ipv6` - Customer peer IPv6 address.
* `direct_connect_interface_id` - The ID of the direct connect interface.
* `ha_direct_connect_interface_id` - Ha direct connect interface ID.
* `ha_direct_connect_interface_name` - Ha direct connect interface name.
* `ha_vlan_id` - The id of vlan in direct connect.
* `local_peer_ipv6` - Local peer IPv6 address.
* `priority` - Priority of the direct connect interface.
* `state` - State.


## Import

DCInterface can be imported using the id, e.g.

```
$ terraform import ksyun_direct_connect_interface.test 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

