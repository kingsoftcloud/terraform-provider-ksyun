---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_routes"
sidebar_current: "docs-ksyun-datasource-routes"
description: |-
  This data source provides a list of Route resources according to their Route ID, cidr and the VPC they belong to.
---

# ksyun_routes

This data source provides a list of Route resources according to their Route ID, cidr and the VPC they belong to.

#

## Example Usage

```hcl
data "ksyun_routes" "default" {
  output_file  = "output_result"
  ids          = []
  vpc_ids      = []
  instance_ids = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Route IDs, all the Route resources belong to this region will be retrieved if the ID is `""`.
* `instance_ids` - (Optional) A list of the Route target id.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `vpc_ids` - (Optional) A list of VPC id that the desired Route belongs to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `routes` - An information list of routes. Each element contains the following attributes:
  * `create_time` - creation time of the route.
  * `destination_cidr_block` - The cidr block of the desired Route.
  * `id` - The ID of Route.
  * `next_hop_set` - A list of next hop.
    * `gateway_id` - ID of the gateway.
    * `gateway_name` - Name of the gateway.
  * `route_id` - The ID of Route.
  * `route_type` - The type of the desired Route.
  * `vpc_id` - The ID of VPC.
* `total_count` - Total number of Route resources that satisfy the condition.


