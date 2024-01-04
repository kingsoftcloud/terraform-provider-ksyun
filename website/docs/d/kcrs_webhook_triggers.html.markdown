---
subcategory: "KCR"
layout: "ksyun"
page_title: "ksyun: ksyun_kcrs_webhook_triggers"
sidebar_current: "docs-ksyun-datasource-kcrs_webhook_triggers"
description: |-
  This data source provides a list of webhook trigger resources according to their instance id.
---

# ksyun_kcrs_webhook_triggers

This data source provides a list of webhook trigger resources according to their instance id.

## Example Usage

```hcl
data "ksyun_kcrs_webhook_triggers" "foo" {
  output_file = "kcrs_webhook_triggers_output_result"
  instance_id = "86f14f8c-bf24-42c8-91bd-xxxxxxxx"
  namespace   = "tftest"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Kcrs Instance Id.
* `namespace` - (Required) Kcrs Instance Namespace.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `trigger_id` - (Optional) Webhook Trigger ID, all the Webhook Trigger belong to this namespace of instance will be retrieved if the ID is `""`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total_count` - Total number of kcrs webhook_triggers that satisfy the condition.
* `triggers` - It is a nested type which documented below.
  * `create_time` - name of the certificate.
  * `enabled` - name of the certificate.
  * `event_type` - name of the certificate.
  * `trigger_id` - ID of the certificate.
  * `trigger_name` - name of the certificate.
  * `trigger_url` - name of the certificate.
  * `update_time` - name of the certificate.


