---
subcategory: "KNAD"
layout: "ksyun"
page_title: "ksyun: ksyun_knad"
sidebar_current: "docs-ksyun-resource-knad"
description: |-
  Provides an KNAD resource.
---

# ksyun_knad

Provides an KNAD resource.

## Example Usage

```hcl
# Create an knad
resource "ksyun_knad" "default" {
  link_type  = "DDoS_BGP"
  ip_count   = 10
  band       = 30
  max_band   = 30
  idc_band   = 100
  duration   = 1
  knad_name  = "ksc_knad"
  bill_type  = 1
  service_id = "KNAD_30G"
  project_id = "0"
}
```

## Argument Reference

The following arguments are supported:

* `band` - (Required) the band of the Knad.
* `bill_type` - (Required, ForceNew) the bill type of the Knad. Valid Values: 1:(PrePaidByMonth),5:(DailyPaidByTransfer).
* `idc_band` - (Required) the idcband of the Knad.
* `ip_count` - (Required) the max ip count that can bind to the Knad,value range: [10, 100].
* `max_band` - (Required) the max band of the Knad.
* `service_id` - (Required) The service id of the Knad.Valid Values:'KNAD_30G','KNAD_100G','KNAD_300G','KNAD_1000G',''KNAD_2000G''.
* `duration` - (Optional, ForceNew) Purchase time.If bill_type is 1,this is Required.
* `knad_name` - (Optional) the name of the Knad.
* `link_type` - (Optional) the link type of the Knad. Valid Values: 'DDoS_BGP'.
* `project_id` - (Optional) The id of the project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `knad_id` - the ID of the Knad.


## Import

Knad can be imported using the id, e.g.

```
$ terraform import ksyun_knad.default knad67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

