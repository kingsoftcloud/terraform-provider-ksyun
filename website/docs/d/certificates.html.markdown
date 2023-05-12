---
subcategory: "KCM"
layout: "ksyun"
page_title: "ksyun: ksyun_certificates"
sidebar_current: "docs-ksyun-datasource-certificates"
description: |-
  This data source provides a list of Certificate resources (KCM) according to their ID.
---

# ksyun_certificates

This data source provides a list of Certificate resources (KCM) according to their ID.

## Example Usage

```hcl
data "ksyun_certificates" "default" {
  output_file = "output_result"
  ids         = ["c7b2ba05-9302-4933-8588-a66f920ff57d"]
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Certificate IDs, all the Certificates belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter results by certificate name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `certificates` - It is a nested type which documented below.
  * `certificate_id` - ID of the certificate.
  * `certificate_name` - name of the certificate.
* `total_count` - Total number of certificates that satisfy the condition.


