---
subcategory: "IAM"
layout: "ksyun"
page_title: "ksyun: ksyun_iam_role"
sidebar_current: "docs-ksyun-resource-iam_role"
description: |-
  Provides a Iam Role resource.
---

# ksyun_iam_role

Provides a Iam Role resource.

#

## Example Usage

```hcl
resource "ksyun_iam_role" "role" {
  role_name      = "role_name_test"
  trust_accounts = "2000096256"
  description    = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `role_name` - (Required, ForceNew) IAM RoleName.
* `trust_accounts` - (Required, ForceNew) IAM TrustAccounts.
* `description` - (Optional, ForceNew) IAM Description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

IAM Role can be imported using the `role_name`, e.g.

```
$ terraform import ksyun_iam_role.role role_name
```

