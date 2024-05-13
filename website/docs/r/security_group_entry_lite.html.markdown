---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_security_group_entry_lite"
sidebar_current: "docs-ksyun-resource-security_group_entry_lite"
description: |-
  Provides a Security Group Entry resource that can manage a list of diverse cidr_block.
---

# ksyun_security_group_entry_lite

Provides a Security Group Entry resource that can manage a list of diverse cidr_block.

#

## Example Usage

```hcl
resource "ksyun_security_group_entry_lite" "default" {
  security_group_id = "7385c8ea-xxxx-xxxx-xxxx-517fc3726256"
  cidr_block        = ["10.0.0.1/32", "10.111.222.1/32"]
  direction         = "in"
  protocol          = "ip"
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required, ForceNew) The cidr block list of security group rule.
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
* `security_group_entry_id_list` - The security group entry id of this lite managed.


## Import

-> **NOTE:** This resource cannot be imported. if you need import security group entry, you are supposed to use `ksyun_security_group_entry`

