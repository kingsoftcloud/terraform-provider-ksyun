---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_network_acls"
sidebar_current: "docs-ksyun-datasource-network_acls"
description: |-
  This data source provides a list of Network ACL resources
---

# ksyun_network_acls

This data source provides a list of Network ACL resources

#

## Example Usage

```hcl
data "ksyun_network_acls" "default" {
  output_file = "output_result"

  //  vpc_ids = ["769c780b-acbd-41ca-9a06-4960e2423c7e"]
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of network ACL IDs.
* `name_regex` - (Optional) A regex string to filter results by ACL name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `vpc_ids` - (Optional) A list of VPC IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `network_acls` - An information list of ACLs. Each element contains the following attributes:
  * `create_time` - creation time of the ACL.
  * `id` - ID of the ACL.
  * `name` - Name of the ACL.
  * `network_acl_entry_set` - A list of the ACL entries.
    * `cidr_block` - The information of Acl's cidr block.
    * `description` - Description of the ACL entry.
    * `direction` - rule direction.
    * `icmp_code` - ICMP code.
    * `icmp_type` - ICMP type.
    * `network_acl_entry_id` - ID of the ACL entry.
    * `network_acl_id` - ID of the ACL.
    * `port_range_from` - beginning of the port range.
    * `port_range_to` - ending of the port range.
    * `protocol` - rule protocol.
    * `rule_action` - rule action, allow or deny.
    * `rule_number` - rule priority.
  * `network_acl_id` - ID of the ACL.
  * `network_acl_name` - Name of the ACL.
  * `vpc_id` - ID of the VPC.
* `total_count` - Total number of ACLs that satisfy the condition.


