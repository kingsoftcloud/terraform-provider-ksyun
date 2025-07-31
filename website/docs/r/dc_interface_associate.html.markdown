---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_dc_interface_associate"
sidebar_current: "docs-ksyun-resource-dc_interface_associate"
description: |-
  Associate an Direct connect gateway resource with interface.
---

# ksyun_dc_interface_associate

Associate an Direct connect gateway resource with interface.

#

## Example Usage

```hcl
resource "ksyun_dc_interface_associate" "test" {
  direct_connect_interface_id = ksyun_direct_connect_interface.test.id
  direct_connect_gateway_id   = ksyun_direct_connect_gateway.test.id
}
```

## Argument Reference

The following arguments are supported:

* `direct_connect_gateway_id` - (Required, ForceNew) The Gateway ID of the direct connect gateway.
* `direct_connect_interface_id` - (Required, ForceNew) The ID of the direct connect interface to be associated with the gateway.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



