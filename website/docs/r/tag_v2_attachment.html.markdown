---
subcategory: "Tag"
layout: "ksyun"
page_title: "ksyun: ksyun_tag_v2_attachment"
sidebar_current: "docs-ksyun-resource-tag_v2_attachment"
description: |-
  Provides an attachment for pinning tag upon resource.
---

# ksyun_tag_v2_attachment

Provides an attachment for pinning tag upon resource.

> Note: supported all of resource_type
> The tag will be created if it is not existed.

#

## Example Usage

```hcl
resource "ksyun_tag_v2" "tagv2" {
  key   = "test_tag_key"
  value = "test_tag_value"
}

resource "ksyun_tag_v2_attachment" "tag" {
  key           = "test_tag_key"
  value         = "test_tag_value"
  resource_type = "redis-instance"
  resource_id   = "1f4e8c22-xxxx-xxxx-xxxx-cc6345011af4"
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required) Tag key.
* `resource_id` - (Required) Resource ID.
* `resource_type` - (Required) Resource type. [supported type](https://docs.ksyun.com/documents/43391).
* `value` - (Required) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `tag_id` - Tag id.


## Import

Tagv2Attachment can be imported using the `id`, e.g.

```
$ terraform import ksyun_tag_v2_attachment.tag ${tag_key}:${tag_value},${resource_type}:${resource_id}
```

