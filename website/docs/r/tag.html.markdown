---
subcategory: "Tag"
layout: "ksyun"
page_title: "ksyun: ksyun_tag"
sidebar_current: "docs-ksyun-resource-tag"
description: |-
  Provides a tag resource.
---

# ksyun_tag

Provides a tag resource.

#

## Example Usage

```hcl
resource "ksyun_tag" "kec_tag" {
  key           = "test_tag_key"
  value         = "test_tag_value"
  resource_type = "eip"
  resource_id   = ' xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx '
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required) Tag key.
* `resource_id` - (Required) Resource ID.
* `resource_type` - (Required) Resource type.
* `value` - (Required) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Tag can be imported using the `id`, e.g.

```
$ terraform import ksyun_tag.kec_tag ${tag_key}:${tag_value},${resource_type}:${resource_id}
```

