---
subcategory: "IAM"
layout: "ksyun"
page_title: "ksyun: ksyun_iam_project"
sidebar_current: "docs-ksyun-resource-iam_project"
description: |-
  Provides a Iam Project resource.
---

# ksyun_iam_project

Provides a Iam Project resource.

#

## Example Usage

```hcl
resource "ksyun_iam_project" "project" {
  project_name = "ProjectNameTest"
  project_desc = "ProjectDesc"
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (Required) IAM ProjectName.
* `project_desc` - (Optional) IAM ProjectDesc.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

IAM Project can be imported using the `project_name`, e.g.

```
$ terraform import ksyun_iam_project.project project_name
```

