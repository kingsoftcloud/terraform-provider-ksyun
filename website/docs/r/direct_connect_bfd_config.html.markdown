---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_direct_connect_bfd_config"
sidebar_current: "docs-ksyun-resource-direct_connect_bfd_config"
description: |-
  Provides a DirectConnectBfdConfig resource.
---

# ksyun_direct_connect_bfd_config

Provides a DirectConnectBfdConfig resource.

## Example Usage

```hcl
resource "ksyun_direct_connect_bfd_config" "test" {
  min_tx_interval   = 100
  min_rx_interval   = 200
  detect_multiplier = 3
  multi_hop         = true
}
```

## Argument Reference

The following arguments are supported:

* `detect_multiplier` - (Optional) Detect Multiplier.
* `min_rx_interval` - (Optional) The interval at which the BFD control packets are received.
* `min_tx_interval` - (Optional) The interval at which the BFD control packets are sent.
* `multi_hop` - (Optional) Whether the multi hop.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `bfd_config_id` - The ID of the BFD configuration.


## Import

ksyun_direct_connect_bfd_config can be imported using the id, e.g.

```
$ terraform import ksyun_direct_connect_bfd_config.test 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

