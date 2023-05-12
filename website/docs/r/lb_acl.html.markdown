---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_acl"
sidebar_current: "docs-ksyun-resource-lb_acl"
description: |-
  Provides a Load Balancer acl resource to add content forwarding policies for Load Balancer backend resource.
---

# ksyun_lb_acl

Provides a Load Balancer acl resource to add content forwarding policies for Load Balancer backend resource.

#

## Example Usage

```hcl
resource "ksyun_lb_acl" "default" {
  load_balancer_acl_name = "tf-xun2"
}
```

## Argument Reference

The following arguments are supported:

* `ip_version` - (Optional) IP version of the load balancer acl. valid values:'ipv4', 'ipv6'. default is 'ipv4'.
* `load_balancer_acl_entry_set` - (Optional) ACL Entries. this parameter will be deprecated, use `ksyun_lb_acl_entry` instead.
* `load_balancer_acl_name` - (Optional) The name of the load balancer acl.

The `load_balancer_acl_entry_set` object supports the following:

* `cidr_block` - (Required) The information of the load balancer Acl's cidr block.
* `protocol` - (Optional) protocol.Valid Values:'ip'.
* `rule_action` - (Optional) The action of load balancer Acl rule. Valid Values:'allow', 'deny'. Default is 'allow'.
* `rule_number` - (Optional) The information of the load balancer Acl's rule priority. value range:[1-32766].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - creation time of the load balancer acl.
* `load_balancer_acl_id` - ID of the load balancer acl.


## Import

LB ACL can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb_acl.example fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```

