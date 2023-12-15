---
subcategory: "Tag"
layout: "ksyun"
page_title: "ksyun: ksyun_tag_v2"
sidebar_current: "docs-ksyun-resource-tag_v2"
description: |-
  Provides a Tagv2 resource.
---

# ksyun_tag_v2

Provides a Tagv2 resource.

#

## Example Usage

```hcl
resource "ksyun_tag_v2" "tag" {
  key   = "test_tag_key"
  value = "test_tag_value"
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required) Tag key.
* `value` - (Required) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Tagv2 can be imported using the `key&value`, e.g.

```
$ terraform import ksyun_tag_v2.tag ${tagv2_key}:${tagv2_value}
```

