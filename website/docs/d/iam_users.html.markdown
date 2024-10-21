---
subcategory: "IAM"
layout: "ksyun"
page_title: "ksyun: ksyun_iam_users"
sidebar_current: "docs-ksyun-datasource-iam_users"
description: |-
  This data source provides a list of user resources.
---

# ksyun_iam_users

This data source provides a list of user resources.

#

## Example Usage

```hcl
data "ksyun_iam_users" "users" {
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `users` - a list of users.
  * `country_mobile_code` - IAN User CountryMobileCode.
  * `create_date` - IAN User CreateDate.
  * `email_verified` - IAN User EmailVerified.
  * `email` - IAN User Email.
  * `enable_mfa` - IAN User EnableMFA.
  * `id` - The ID of the IAM User Id.
  * `is_international` - IAN User IsInternational.
  * `krn` - IAN User Krn.
  * `password_reset_required` - IAN User PasswordResetRequired.
  * `path` - IAM User Path.
  * `phone_verified` - IAN User PhoneVerified.
  * `phone` - IAN User Phone.
  * `real_name` - IAM User RealName.
  * `remark` - IAN User Remark.
  * `update_date` - IAN User UpdateDate.
  * `user_id` - The ID of the IAM UserId.
  * `user_name` - IAM UserName.


