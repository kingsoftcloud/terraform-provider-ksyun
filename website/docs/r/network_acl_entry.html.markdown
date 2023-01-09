---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_network_acl_entry"
sidebar_current: "docs-ksyun-resource-network_acl_entry"
description: |-
  Provides a Network ACL Entry resource under Network ACL resource.
---

# ksyun_network_acl_entry

Provides a Network ACL Entry resource under Network ACL resource.

#

## Example Usage

```hcl
resource "ksyun_network_acl_entry" "test" {
  description    = "测试1"
  cidr_block     = "10.0.16.0/24"
  rule_number    = 16
  direction      = "in"
  rule_action    = "deny"
  protocol       = "ip"
  network_acl_id = "679b6a88-67dd-4e17-a80a-985d9673050e"
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required, ForceNew) The cidr_block of the network acl entry.
* `direction` - (Required, ForceNew) The direction of the network acl entry. Valid Values: 'in','out'.
* `network_acl_id` - (Required, ForceNew) The id of the network acl.
* `protocol` - (Required, ForceNew) The protocol of the network acl entry.Valid Values: 'ip','icmp','tcp','udp'.
* `rule_action` - (Required, ForceNew) The rule_action of the network acl entry.Valid Values: 'allow','deny'.
* `rule_number` - (Required, ForceNew) The rule_number of the network acl entry. value range:[1,32766].
* `description` - (Optional) The description of the network acl entry.
* `icmp_code` - (Optional, ForceNew) The icmp_code of the network acl entry.If protocol is icmp, Required.
* `icmp_type` - (Optional, ForceNew) The icmp_type of the network acl entry.If protocol is icmp, Required.
* `port_range_from` - (Optional, ForceNew) The port_range_from of the network acl entry.If protocol is tcp or udp,Required.
* `port_range_to` - (Optional, ForceNew) The port_range_to of the network acl entry.If protocol is tcp or udp,Required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `network_acl_entry_id` - ID of the network acl entry.


