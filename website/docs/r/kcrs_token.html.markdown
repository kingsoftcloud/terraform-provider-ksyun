---
subcategory: "KCR"
layout: "ksyun"
page_title: "ksyun: ksyun_kcrs_token"
sidebar_current: "docs-ksyun-resource-kcrs_token"
description: |-
  Provides an access token resource under kcrs repository instance.
---

# ksyun_kcrs_token

Provides an access token resource under kcrs repository instance.

## Example Usage

```hcl
# repository instance
resource "ksyun_kcrs_instance" "foo" {
  instance_name = "tfunittest"
  instance_type = "basic"
}

# Create a KcrsToken
resource "ksyun_kcrs_token" "foo" {
  instance_id = ksyun_kcrs_instance.foo.id
  token_type  = "Day"
  token_time  = 10
  desc        = "test"
  enable      = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Instance id of repository.
* `token_time` - (Required) The validation time of token. If the `token_type` is 'NeverExpire', this field is invalid.
* `token_type` - (Required) Token type.
* `desc` - (Optional) Description for this token.
* `enable` - (Optional) Whether to enable this token.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `expire_time` - The expired time for this token.


## Import

KcrsToken can be imported using `instance_id:token_id`, e.g.

```
$ terraform import ksyun_kcrs_token.foo ${instance_id}:${token_id}
```

