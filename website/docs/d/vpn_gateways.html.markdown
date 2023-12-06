---
subcategory: "VPN"
layout: "ksyun"
page_title: "ksyun: ksyun_vpn_gateways"
sidebar_current: "docs-ksyun-datasource-vpn_gateways"
description: |-
  This data source provides a list of VPN gateways.
---

# ksyun_vpn_gateways

This data source provides a list of VPN gateways.

#

## Example Usage

```hcl
data "ksyun_vpn_gateways" "default" {
  output_file = "output_result"
  ids         = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of VPN gateway IDs, all the resources belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter results by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `project_ids` - (Optional) A list of project IDs.
* `vpc_ids` - (Optional) A list of VPC IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total_count` - Total number of resources that satisfy the condition.
* `vpn_gateways` - It is a nested type which documented below.
  * `band_width` - Band width.
  * `create_time` - Creation time.
  * `gateway_address` - Gateway IP address.
  * `ha_gateway_address` - HA Gateway IP address.
  * `id` - The ID of the gateway.
  * `name` - The name of the gateway.
  * `remote_cidr_set` - A list of remote cidrs.
    * `cidr_block` - Cidr block.
  * `vpc_id` - VPC ID.
  * `vpn_gateway_id` - The ID of the gateway.
  * `vpn_gateway_name` - The name of the gateway.
  * `vpn_gateway_version` - The version of vpn gateway.


