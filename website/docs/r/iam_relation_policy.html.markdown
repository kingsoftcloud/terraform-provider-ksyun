---
subcategory: "IAM"
layout: "ksyun"
page_title: "ksyun: ksyun_iam_relation_policy"
sidebar_current: "docs-ksyun-resource-iam_relation_policy"
description: |-
  Provides a Iam Policy resource.
---

# ksyun_iam_relation_policy

Provides a Iam Policy resource.

#

## Example Usage

```hcl
resource "ksyun_iam_relation_policy" "user" {
  name          = "iam_user_name"
  policy_name   = "IAMReadOnlyAccess"
  relation_type = 1
  policy_type   = "system"
} `

resource "ksyun_iam_relation_policy" "user" {
  name          = "iam_role_name"
  policy_name   = "IAMReadOnlyAccess"
  relation_type = 2
  policy_type   = "system"
} `
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) IAM UserName or RoleName according to relation type.
* `policy_name` - (Required, ForceNew) IAM PolicyName.
* `policy_type` - (Required, ForceNew) policy type system is the system policy,policy type custom is the custom policy.
* `relation_type` - (Required, ForceNew) relation type 1 is the user,relation type 2 is the role.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

IAM Policy can be imported using the `policy_name`, e.g.

```
$ terraform import ksyun_iam_relation_policy.user
```

