---
subcategory: "IAM"
layout: "ksyun"
page_title: "ksyun: ksyun_iam_groups"
sidebar_current: "docs-ksyun-datasource-iam_groups"
description: |-
  This data source provides a list of group resources.
---

# ksyun_iam_groups

This data source provides a list of group resources.

#

## Example Usage

```hcl
data "ksyun_iam_groups" "groups" {
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `groups` - a list of users.
  * `create_date` - IAN Group CreateDate.
  * `description` - IAM Group Description.
  * `group_id` - The ID of the IAM GroupId.
  * `group_name` - IAM GroupName.
  * `krn` - IAN Group Krn.
  * `path` - IAM Group Path.
  * `policy_count` - IAN Group PolicyCount.
  * `user_count` - IAN Group UserCount.


