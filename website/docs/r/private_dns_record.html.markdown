---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_private_dns_record"
sidebar_current: "docs-ksyun-resource-private_dns_record"
description: |-
  Provides a Private Dns Record resource under PDNS Zone resource.
---

# ksyun_private_dns_record

Provides a Private Dns Record resource under PDNS Zone resource.

#

## Example Usage

```hcl
resource "ksyun_private_dns_zone" "foo" {
  zone_name   = "tf-pdns-zone-pdns.com"
  zone_ttl    = 360
  charge_type = "TrafficMonthly"
}

resource "ksyun_private_dns_record" "foo" {
  record_name  = "tf-pdns-record"
  record_ttl   = 360
  zone_id      = ksyun_private_dns_zone.foo.id
  type         = "CNAME"
  record_value = "tf-record.com"
}
```

## Argument Reference

The following arguments are supported:

* `record_name` - (Required, ForceNew) The record name of private dns.
* `record_value` - (Required, ForceNew) Record value, such as IP: 192.168.10.2, CNAME: cname.ksyun.com, and MX: mail.ksyun.com..
* `type` - (Required, ForceNew) Record type. Valid values: "A", "AAAA", "CNAME", "MX", "TXT", "SRV".
* `zone_id` - (Required, ForceNew) Private Dns Zone ID.
* `port` - (Optional, ForceNew) The port of record in which is associated with domain or ip. Required, when type is `SRV`.
* `priority` - (Optional, ForceNew) Record priority. Value range: [SRV|0~65535], [MX|1~100]. Required, when type is `SRV` or `MX`.
* `record_ttl` - (Optional) Record cache time. The smaller the value, the faster the record will take effect. Value range: 60~86400s.
* `weight` - (Optional, ForceNew) Record weight. Value range: 0~65535. Required, when type is `SRV`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



