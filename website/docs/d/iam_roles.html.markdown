---
subcategory: "IAM"
layout: "ksyun"
page_title: "ksyun: ksyun_iam_roles"
sidebar_current: "docs-ksyun-datasource-iam_roles"
description: |-
  This data source provides a list of role resources.
---

# ksyun_iam_roles

This data source provides a list of role resources.

#

## Example Usage

```hcl
data "ksyun_iam_roles" "roles" {
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `roles` - a list of users.
  * `create_date` - IAN Role CreateDate.
  * `description` - IAM Role Description.
  * `krn` - IAM Role Krn.
  * `role_id` - The ID of the IAM RoleId.
  * `role_name` - IAM RoleName.
  * `service_role_type` - IAN Role ServiceRoleType.
  * `trust_accounts` - IAN Role TrustAccounts.
  * `trust_provider` - IAN Role TrustProvider.
  * `trust_type` - IAN Role TrustType.


