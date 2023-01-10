---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_security_group"
sidebar_current: "docs-ksyun-resource-security_group"
description: |-
  Provides a Security Group resource.
---

# ksyun_security_group

Provides a Security Group resource.

#

## Example Usage

```hcl
resource "ksyun_security_group" "default" {
  vpc_id              = "26231a41-4c6b-4a10-94ed-27088d5679df"
  security_group_name = "xuan-tf--s"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The Id of the vpc.
* `security_group_entries` - (Optional) Network security group Entries. this parameter will be deprecated, use `ksyun_security_group_entry` instead.
* `security_group_name` - (Optional) The name of the security group.

The `security_group_entries` object supports the following:

* `cidr_block` - (Required) The cidr block of security group rule.
* `direction` - (Required) The direction of the entry, valid values:'in', 'out'.
* `protocol` - (Required) The protocol of the entry, valid values: 'ip', 'tcp', 'udp', 'icmp'.
* `description` - (Optional) The description of the entry.
* `icmp_code` - (Optional) ICMP code.The required if protocol type is 'icmp'.
* `icmp_type` - (Optional) ICMP type.The required if protocol type is 'icmp'.
* `port_range_from` - (Optional) Port rule start port for TCP or UDP protocol.The required if protocol type is 'tcp' or 'udp'.
* `port_range_to` - (Optional) Port rule end port for TCP or UDP protocol.The required if protocol type is 'tcp' or 'udp'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time of creation of security group.
* `security_group_id` - The ID of the security group.


## Import

Security Group can be imported using the `id`, e.g.

```
$ terraform import ksyun_security_group.example xxxxxxxx-abc123456
```

