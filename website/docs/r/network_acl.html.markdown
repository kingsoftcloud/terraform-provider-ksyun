---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_network_acl"
sidebar_current: "docs-ksyun-resource-network_acl"
description: |-
  Provides a Network ACL resource under VPC resource.
---

# ksyun_network_acl

Provides a Network ACL resource under VPC resource.

#

## Example Usage

```hcl
resource "ksyun_network_acl" "default" {
  vpc_id           = "a8979fe2-cf1a-47b9-80f6-57445227c541"
  network_acl_name = "ceshi"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The id of the vpc.
* `network_acl_entries` - (Optional) Network ACL Entries. this parameter will be deprecated, use `ksyun_network_acl_entry` instead.
* `network_acl_name` - (Optional) The name of the network ACL.

The `network_acl_entries` object supports the following:

* `cidr_block` - (Required) The cidr_block of the network acl entry.
* `direction` - (Required) The direction of the network acl entry. Valid Values: 'in','out'.
* `protocol` - (Required) The protocol of the network acl entry.Valid Values: 'ip','icmp','tcp','udp'.
* `rule_action` - (Required) The rule_action of the network acl entry.Valid Values: 'allow','deny'.
* `rule_number` - (Required) The rule_number of the network acl entry. value range:[1,32766].
* `description` - (Optional) The description of the network acl entry.
* `icmp_code` - (Optional) The icmp_code of the network acl entry.If protocol is icmp, Required.
* `icmp_type` - (Optional) The icmp_type of the network acl entry.If protocol is icmp, Required.
* `port_range_from` - (Optional) The port_range_from of the network acl entry.If protocol is tcp or udp,Required.
* `port_range_to` - (Optional) The port_range_to of the network acl entry.If protocol is tcp or udp,Required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The creation time of network acl.
* `network_acl_id` - The ID of the network ACL.


## Import

Network ACL can be imported using the `id`, e.g.

```
$ terraform import ksyun_network_acl.default fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```

