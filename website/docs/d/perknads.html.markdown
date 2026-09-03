---
subcategory: "KNAD"
layout: "ksyun"
page_title: "ksyun: ksyun_perknads"
sidebar_current: "docs-ksyun-datasource-perknads"
description: |-
  This data source provides a list of PerPay KNAD resources according to their KNAD ID.
---

# ksyun_perknads

This data source provides a list of PerPay KNAD resources according to their KNAD ID.

#

## Example Usage

```hcl
data "ksyun_perknads" "default" {
  output_file = "output_result"

  project_id = []
  ids        = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Knad IDs, all the PerKnads belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `project_id` - (Optional) One or more project IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `perknads` - An information list of PerPay Knad. Each element contains the following attributes:
  * `band` - the band of the Knad.
  * `bill_type` - the bill type of the Knad.
  * `exprie_time` - the exprie time of the Knad.
  * `ip_count` - the max ip count that can bind to the Knad.
  * `knad_id` - ID of the resource.
  * `knad_name` - the name of the Knad.
  * `max_band` - the max band of the Knad.
  * `project_id` - The id of the project.
  * `used_ip_count` - The binding ip count of the Knad.


