---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_direct_connect_gateway"
sidebar_current: "docs-ksyun-resource-direct_connect_gateway"
description: |-
  Provides a DirectConnectGateWay resource.
---

# ksyun_direct_connect_gateway

Provides a DirectConnectGateWay resource.

## Example Usage

```hcl
resource "ksyun_direct_connect_gateway" "test" {
  direct_connect_gateway_name = "tf_direct_connect_gateway_test_1"
  vpc_id                      = "a38673ae-c9b7-4f8e-b727-b6feb648xxxx"
}
```

## Argument Reference

The following arguments are supported:

* `direct_connect_gateway_name` - (Optional) The name of the direct connect gateway.
* `vpc_id` - (Optional) Vpc Id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `associated_instance_type` - The type of associated instance.
* `band_width` - Band width.
* `cen_account_id` - The id of the CEN account associated with the direct connect gateway.
* `cen_id` - ID of the Cen.
* `direct_connect_gateway_id` - The ID of the direct connect gateway.
* `direct_connect_interface_info_set` - The set of direct connect associated interface info.
  * `direct_connect_interface_id` - ID of the direct connect interface.
* `direct_connect_interface_name` - The name of the direct connect interface associated with the direct connect gateway.
* `nat_id` - The ID of the NAT gateway associated with the direct connect gateway.
* `remote_cidr_set` - The set of remote cidr.
* `status` - Status.
* `version` - Version.


## Import

ksyun_direct_connect_gateway can be imported using the id, e.g.

```
$ terraform import ksyun_direct_connect_gateway.test 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

