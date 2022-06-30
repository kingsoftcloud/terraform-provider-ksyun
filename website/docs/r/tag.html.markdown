---
layout: "ksyun"
page_title: "Ksyun: ksyun_tag"
sidebar_current: "docs-ksyun-resource-tag"
description: |-
  Provides a tag resource.
---

# ksyun_tag

Provides a tag resource.

## Example Usage

```hcl
resource "ksyun_tag" "kec_tag" {
  key = "test_tag_key"
  value = "test_tag_value"
  resource_type = "eip"
  resource_id = 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required) The key of tag.
* `value` - (Required) The value of tag.
* `resource_type` - (Required) The type of the instance. Valid Values: kec-instance、rds-instance、kcs-instance、epc-instance、eip、ebs、slb、nat、bws、peering、ks3
* `resource_id` - (Required) The id of the interface.

