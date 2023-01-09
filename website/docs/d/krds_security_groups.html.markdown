---
subcategory: "KRDS"
layout: "ksyun"
page_title: "ksyun: ksyun_krds_security_groups"
sidebar_current: "docs-ksyun-datasource-krds_security_groups"
description: |-
  Query security group information
---

# ksyun_krds_security_groups

Query security group information

#

## Example Usage

```hcl
# Get  krds_security_groups

data "ksyun_krds_security_groups" "security_groups" {
  output_file       = "output_file"
  security_group_id = 123
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Required) The filename of the content store will be returned.
* `security_group_id` - (Optional) Security group ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_groups` - An information list of KRDS security groups.
  * `created` - The time of creation.
  * `instances` - corresponding instance.
    * `created` - The time of creation.
    * `db_instance_identifier` - instance ID.
    * `db_instance_name` - instance name.
    * `db_instance_type` - instance type.
    * `vip` - instance virtual IP.
  * `security_group_description` - Security group description.
  * `security_group_id` - Security group ID.
  * `security_group_name` - Security group name.
  * `security_group_rules` - security group rules.
    * `created` - The time of creation.
    * `security_group_rule_id` - rule ID.
    * `security_group_rule_name` - rule name.
    * `security_group_rule_protocol` - rule protocol.
* `total_count` - Total number of resources that satisfy the condition.


