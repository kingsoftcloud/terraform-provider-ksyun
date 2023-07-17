---
subcategory: "KNAD"
layout: "ksyun"
page_title: "ksyun: ksyun_knads"
sidebar_current: "docs-ksyun-datasource-knads"
description: |-
  This data source provides a list of KNAD resources  according to their KNAD ID.
---

# ksyun_knads

This data source provides a list of KNAD resources  according to their KNAD ID.

#

## Example Usage

```hcl
data "ksyun_knads" "default" {
  output_file = "output_result"

  project_id = []
  ids        = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Knad IDs, all the Knads belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `project_id` - (Optional) One or more project IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `knads` - An information list of Knad. Each element contains the following attributes:
  * `band` - the band of the Knad.
  * `bill_type` - the bill type of the Knad. Valid Values: 1:(PrePaidByMonth),5:(DailyPaidByTransfer).
  * `exprie_time` - the exprie time of the Knad.
  * `idc_band` - the idcband of the Knad.
  * `ip_count` - the max ip count that can bind to the Knad.
  * `knad_id` - ID of the resource.
  * `knad_name` - the name of the Knad.
  * `max_band` - the max band of the Knad.
  * `project_id` - The id of the project.
  * `service_id` - the service id of the Knad.Valid Values:'KNAD_30G','KNAD_100G','KNAD_300G','KNAD_1000G',''KNAD_2000G''.
  * `used_ip_count` - The binding ip count of the Knad.


