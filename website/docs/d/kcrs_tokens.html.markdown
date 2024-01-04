---
subcategory: "KCR"
layout: "ksyun"
page_title: "ksyun: ksyun_kcrs_tokens"
sidebar_current: "docs-ksyun-datasource-kcrs_tokens"
description: |-
  This data source provides a list of token resources according to their instance id.
---

# ksyun_kcrs_tokens

This data source provides a list of token resources according to their instance id.

## Example Usage

```hcl
data "ksyun_kcrs_tokens" "foo" {
  output_file = "kcrs_tokens_output_result"
  instance_id = "86f14f8c-bf24-42c8-91bd-xxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Kcrs Instance Id.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tokens` - It is a nested type which documented below.
  * `create_time` - Created Time.
  * `desc` - Description for this token.
  * `enable` - Whether Enable.
  * `expire_time` - Expired Time.
  * `token_id` - ID of the token.
* `total_count` - Total number of kcrs tokens that satisfy the condition.


