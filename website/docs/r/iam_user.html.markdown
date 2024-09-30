---
subcategory: "IAM"
layout: "ksyun"
page_title: "ksyun: ksyun_iam_user"
sidebar_current: "docs-ksyun-resource-iam_user"
description: |-
  Provides a Iam User resource.
---

# ksyun_iam_user

Provides a Iam User resource.

#

## Example Usage

```hcl
resource "ksyun_iam_user" "user" {
  user_name                = "username01"
  real_name                = "realname01"
  phone                    = "13800000000"
  email                    = "test@ksyun.com"
  remark                   = "remark"
  password                 = "password"
  password_reset_required  = 0
  open_login_protection    = 1
  open_security_protection = 1
  view_all_project         = 0
}
```

## Argument Reference

The following arguments are supported:

* `user_name` - (Required) IAM UserName.
* `email` - (Optional) IAM Email.
* `open_login_protection` - (Optional) Does IAM user enable login protection.
* `open_security_protection` - (Optional) Does IAM user enable operation protection.
* `password_reset_required` - (Optional) Does IAM user login reset password.
* `password` - (Optional) IAM Password.
* `phone` - (Optional) IAM Phone.
* `real_name` - (Optional) IAM RealName.
* `remark` - (Optional) IAM Remark.
* `view_all_project` - (Optional) Can IAM users view all projects.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

IAM User can be imported using the `user_name`, e.g.

```
$ terraform import ksyun_iam_user.user user_name
```

