---
subcategory: "KRDS"
layout: "ksyun"
page_title: "ksyun: ksyun_krds_security_group"
sidebar_current: "docs-ksyun-resource-krds_security_group"
description: |-
  Provide RDS security group function
---

# ksyun_krds_security_group

Provide RDS security group function

#

## Example Usage

```hcl
resource "ksyun_krds_security_group" "default" {
  security_group_name        = "terraform_security_group_13"
  security_group_description = "terraform-security-group-13"
  security_group_rule {
    security_group_rule_protocol = "182.133.0.0/16"
    security_group_rule_name     = "asdf"
  }
  security_group_rule {
    security_group_rule_protocol = "182.134.0.0/16"
    security_group_rule_name     = "asdf2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `security_group_description` - (Optional) description of security group.
* `security_group_name` - (Optional) the name of the security group.
* `security_group_rule` - (Optional) the rule.

The `security_group_rule` object supports the following:

* `security_group_rule_protocol` - (Required) the protocol of the rule.
* `security_group_rule_name` - (Optional) the name of the rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created` - the creation time of the resource.
* `security_group_id` - Security group ID.


## Import

RDS security group can be imported using the id, e.g.

```
$ terraform import ksyun_krds_security_group.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

