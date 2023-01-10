---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_security_group_entry"
sidebar_current: "docs-ksyun-resource-security_group_entry"
description: |-
  Provides a Security Group Entry resource.
---

# ksyun_security_group_entry

Provides a Security Group Entry resource.

#

## Example Usage

```hcl
resource "ksyun_security_group_entry" "default" {
  security_group_id = "7385c8ea-79f7-4e9c-b99f-517fc3726256"
  cidr_block        = "10.0.0.1/32"
  direction         = "in"
  protocol          = "ip"
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required, ForceNew) The cidr block of security group rule.
* `direction` - (Required, ForceNew) The direction of the entry, valid values:'in', 'out'.
* `protocol` - (Required, ForceNew) The protocol of the entry, valid values: 'ip', 'tcp', 'udp', 'icmp'.
* `security_group_id` - (Required, ForceNew) The ID of the security group.
* `description` - (Optional) The description of the entry.
* `icmp_code` - (Optional, ForceNew) ICMP code.The required if protocol type is 'icmp'.
* `icmp_type` - (Optional, ForceNew) ICMP type.The required if protocol type is 'icmp'.
* `port_range_from` - (Optional, ForceNew) Port rule start port for TCP or UDP protocol.The required if protocol type is 'tcp' or 'udp'.
* `port_range_to` - (Optional, ForceNew) Port rule end port for TCP or UDP protocol.The required if protocol type is 'tcp' or 'udp'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `security_group_entry_id` - The ID of the entry.


## Import

Security Group Entry can be imported using the `id`, e.g.

```
$ terraform import ksyun_security_group_entry.example xxxxxxxx-abc123456
```

