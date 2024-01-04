---
subcategory: "KCR"
layout: "ksyun"
page_title: "ksyun: ksyun_kcrs_namespaces"
sidebar_current: "docs-ksyun-datasource-kcrs_namespaces"
description: |-
  This data source provides a list of namespace resources according to their instance id.
---

# ksyun_kcrs_namespaces

This data source provides a list of namespace resources according to their instance id.

## Example Usage

```hcl
data "ksyun_kcrs_namespaces" "foo" {
  output_file = "kcrs_namespaces_output_result"
  instance_id = "86f14f8c-bf24-42c8-91bd-xxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Kcrs Instance Id.
* `namespace` - (Optional) Kcrs Instance namespace, all the Kcrs namespace belong to this instance will be retrieved if the namespaces is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `namespace_items` - It is a nested type which documented below.
  * `create_time` - Created Time.
  * `namespace` - Namespace.
  * `public` - Whether Public.
  * `repo_count` - The count of Images in this repository.
* `total_count` - Total number of kcrs namespaces that satisfy the condition.


