---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_private_dns_zone"
sidebar_current: "docs-ksyun-resource-private_dns_zone"
description: |-
  Provides a Private Dns Zone resource.
---

# ksyun_private_dns_zone

Provides a Private Dns Zone resource.

#

## Example Usage

```hcl
resource "ksyun_private_dns_zone" "foo" {
  zone_name   = "tf-pdns-zone-pdns.com"
  zone_ttl    = 360
  charge_type = "TrafficMonthly"
}
```

## Argument Reference

The following arguments are supported:

* `charge_type` - (Required, ForceNew) The charge type of the Private Dns Zone. Values: `TrafficMonthly`.
* `zone_name` - (Required, ForceNew) The zone name of private dns.
* `project_id` - (Optional, ForceNew) ID of the project.
* `zone_ttl` - (Optional) The zone cache time. The smaller the value, the faster the record will take effect. Value range: 60~86400s.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `bind_vpc_set` - This zone have bound VPC set.
  * `region_name` - Region name.
  * `status` - The status of binding VPC.
  * `vpc_id` - VPC ID.
  * `vpc_name` - The VPC name.


## Import

Private Dns Record can be imported using the `id`, e.g.

```
$ terraform import ksyun_private_dns_zone.foo fdeba8ca-8aa6-4cd0-8ffa-xxxxxxxxxx
```

