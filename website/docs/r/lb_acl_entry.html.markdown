---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_acl_entry"
sidebar_current: "docs-ksyun-resource-lb_acl_entry"
description: |-
  Provides a Load Balancer acl entry resource to add content forwarding policies for Load Balancer backend resource.
---

# ksyun_lb_acl_entry

Provides a Load Balancer acl entry resource to add content forwarding policies for Load Balancer backend resource.

#

## Example Usage

```hcl
resource "ksyun_lb_acl_entry" "default" {
  load_balancer_acl_id = "8e6d0871-da8a-481e-8bee-b3343e2a6166"
  cidr_block           = "192.168.11.2/32"
  rule_number          = 10
  rule_action          = "allow"
  protocol             = "ip"
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required, ForceNew) The information of the load balancer Acl's cidr block.
* `load_balancer_acl_id` - (Required, ForceNew) The ID of the load balancer acl.
* `protocol` - (Optional, ForceNew) protocol.Valid Values:'ip'.
* `rule_action` - (Optional) The action of load balancer Acl rule. Valid Values:'allow', 'deny'. Default is 'allow'.
* `rule_number` - (Optional) The information of the load balancer Acl's rule priority. value range:[1-32766].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `load_balancer_acl_entry_id` - ID of the LB ACL entry.


## Import

LB ACL entry can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb_acl_entry.example fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```

