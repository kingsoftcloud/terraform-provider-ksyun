---
subcategory: "KCR"
layout: "ksyun"
page_title: "ksyun: ksyun_kcrs_webhook_trigger"
sidebar_current: "docs-ksyun-resource-kcrs_webhook_trigger"
description: |-
  Provides a webhook trigger under kcrs repository instance.
---

# ksyun_kcrs_webhook_trigger

Provides a webhook trigger under kcrs repository instance.

## Example Usage

```hcl
# repository instance
resource "ksyun_kcrs_instance" "foo" {
  instance_name = "tfunittest"
  instance_type = "basic"
}

# Create a webhook trigger
resource "ksyun_kcrs_webhook_trigger" "foo" {
  instance_id = ksyun_kcrs_instance.foo.id
  namespace   = "namespace"
  trigger {
    trigger_url  = "http://www.test111.com"
    trigger_name = "tfunittest"
    event_types  = ["DeleteImage", "PushImage"]
    headers {
      key   = "pp1"
      value = "22"
    }
    headers {
      key   = "pp1"
      value = "333"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, ForceNew) Namespace name.
* `trigger` - (Required) trigger parameters.
* `instance_id` - (Optional) Instance id of repository.

The `headers` object supports the following:

* `key` - (Required) Header Key.
* `value` - (Required) Header Value.

The `trigger` object supports the following:

* `event_types` - (Required) Trigger action. Valid Values: 'PushImage', 'DeleteImage'.
* `trigger_name` - (Required) Trigger name.
* `trigger_url` - (Required) The post url for webhook after trigger action.
* `enabled` - (Optional) Enable trigger.
* `headers` - (Optional) Custom Headers.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



