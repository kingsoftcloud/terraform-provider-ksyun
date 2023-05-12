---
subcategory: "KRDS"
layout: "ksyun"
page_title: "ksyun: ksyun_krds_security_group_rule"
sidebar_current: "docs-ksyun-resource-krds_security_group_rule"
description: |-
  Provide RDS security group rule
---

# ksyun_krds_security_group_rule

Provide RDS security group rule

#

## Example Usage

```hcl
resource "ksyun_krds_security_group_rule" "default" {
  security_group_rule_protocol = "182.133.0.0/16"
  security_group_id            = "62540"
}
```

## Argument Reference

The following arguments are supported:

* `security_group_id` - (Required, ForceNew) security group id.
* `security_group_rule_protocol` - (Required, ForceNew) security group rule protocol.
* `security_group_rule_name` - (Optional, ForceNew) security group rule name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created` - the creation time.
* `security_group_rule_id` - security group rule id.


## Import

RDS security group rule can be imported using the id, e.g.

```
$ terraform import ksyun_krds_security_group_rule.default ${security_group_id}:${security_group_rule_protocol}
```

