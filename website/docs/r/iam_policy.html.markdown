---
subcategory: "IAM"
layout: "ksyun"
page_title: "ksyun: ksyun_iam_policy"
sidebar_current: "docs-ksyun-resource-iam_policy"
description: |-
  Provides a Iam Policy resource.
---

# ksyun_iam_policy

Provides a Iam Policy resource.

#

## Example Usage

```hcl
resource "ksyun_iam_policy" "policy" {
  policy_name     = "TestPolicy1"
  policy_document = "{\"Version\": \"2015-11-01\",\"Statement\": [{\"Effect\": \"Allow\",\"Action\": [\"iam:List*\"],\"Resource\": [\"*\"]}]}"
} `
```

## Argument Reference

The following arguments are supported:

* `policy_name` - (Required, ForceNew) IAM PolicyName.
* `policy_document` - (Optional) IAM PolicyDocument.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `policy_krn` - IAM PolicyKrn.


## Import

IAM Policy can be imported using the `policy_name`, e.g.

```
$ terraform import ksyun_iam_policy.policy policy_name
```

