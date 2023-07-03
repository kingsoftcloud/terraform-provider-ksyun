---
subcategory: "VPN"
layout: "ksyun"
page_title: "ksyun: ksyun_vpn_customer_gateways"
sidebar_current: "docs-ksyun-datasource-vpn_customer_gateways"
description: |-
  This data source provides a list of VPN custom gateways.
---

# ksyun_vpn_customer_gateways

This data source provides a list of VPN custom gateways.

#

## Example Usage

```hcl
data "ksyun_vpn_customer_gateways" "default" {
  output_file = "output_result"
  ids         = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of VPN customer gateway IDs, all the resources belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter results by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `customer_gateways` - It is a nested type which documented below.
  * `create_time` - creation time.
  * `customer_gateway_address` - The IP address of the customer gateway.
  * `customer_gateway_id` - The ID of the customer gateway.
  * `customer_gateway_name` - The name of the customer gateway.
  * `ha_customer_gateway_address` - The IP address of the HA customer gateway.
  * `id` - The ID of the customer gateway.
  * `name` - The name of the customer gateway.
* `total_count` - Total number of resources that satisfy the condition.


