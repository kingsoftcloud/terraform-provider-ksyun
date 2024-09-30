---
subcategory: "IAM"
layout: "ksyun"
page_title: "ksyun: ksyun_iam_projects"
sidebar_current: "docs-ksyun-datasource-iam_projects"
description: |-
  This data source provides a list of project resources.
---

# ksyun_iam_projects

This data source provides a list of project resources.

#

## Example Usage

```hcl
data "ksyun_iam_projects" "projects" {
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `projects` - a list of users.
  * `account_id` - IAM Project AccountId.
  * `create_date` - IAN Role CreateDate.
  * `krn` - IAN Project Krn.
  * `project_desc` - IAM Project ProjectDesc.
  * `project_id` - The ID of the IAM ProjectId.
  * `project_name` - IAM Project ProjectName.
  * `status` - IAN Project Status.


