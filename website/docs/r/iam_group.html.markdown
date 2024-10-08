---
subcategory: "IAM"
layout: "ksyun"
page_title: "ksyun: ksyun_iam_group"
sidebar_current: "docs-ksyun-resource-iam_group"
description: |-
  Provides a Iam Group resource.
---

# ksyun_iam_group

Provides a Iam Group resource.

#

## Example Usage

```hcl
resource "ksyun_iam_group" "group" {
  group_name  = "GroupNameTest"
  description = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required, ForceNew) IAM GroupName.
* `description` - (Optional, ForceNew) IAM Group Description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

IAM Group can be imported using the `group_name`, e.g.

```
$ terraform import ksyun_iam_group.group group_name
```

