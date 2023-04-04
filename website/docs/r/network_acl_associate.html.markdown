---
subcategory: "VPC"
layout: "ksyun"
page_title: "ksyun: ksyun_network_acl_associate"
sidebar_current: "docs-ksyun-resource-network_acl_associate"
description: |-
  Provides a Network ACL Associate resource.
---

# ksyun_network_acl_associate

Provides a Network ACL Associate resource.

#

## Example Usage

```hcl
resource "ksyun_network_acl_associate" "test" {
  network_acl_id = "679b6a88-67dd-4e17-a80a-985d9673050e"
  subnet_id      = "84cc79f3-dc88-4f00-a66a-c7e8d68ec615"
}
```

## Argument Reference

The following arguments are supported:

* `network_acl_id` - (Required, ForceNew) The id of the network acl.
* `subnet_id` - (Required, ForceNew) The id of the Subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Network ACL Associate can be imported using the `network_acl_id:subnet_id`, e.g.

```
$ terraform import ksyun_network_acl_associate.default $network_acl_id:$subnet_id
```

