---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb_acls"
sidebar_current: "docs-ksyun-datasource-lb_acls"
description: |-
  This data source provides a list of Load Balancer Rule resources according to their Load Balancer Rule ID.
---

# ksyun_lb_acls

This data source provides a list of Load Balancer Rule resources according to their Load Balancer Rule ID.

#

## Example Usage

```hcl
data "ksyun_lb_acls" "default" {
  output_file = "output_result"
  ids         = []

}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of LB Rule IDs, all the LB Rules belong to the Load Balancer listener will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `lb_acls` - It is a nested type which documented below.
  * `create_time` - The time of creation for LB ACL.
  * `load_balancer_acl_entry_set` - A list of ACL entries.
    * `cidr_block` - The information of Acl's cidr block.
    * `load_balancer_acl_entry_id` - ID of the ACL entry.
    * `load_balancer_acl_id` - ID of the ACL.
    * `protocol` - rul protocol.
    * `rule_action` - rule action, allow or deny.
    * `rule_number` - rule priority.
  * `load_balancer_acl_id` - ID of the LB ACL.
  * `load_balancer_acl_name` - Name of the LB ACL.
* `total_count` - Total number of LB Rules that satisfy the condition.


