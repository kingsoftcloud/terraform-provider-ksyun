---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_kec_network_interface_attachment"
sidebar_current: "docs-ksyun-resource-kec_network_interface_attachment"
description: |-
  Provides a KEC network interface attachment resource
---

# ksyun_kec_network_interface_attachment

Provides a KEC network interface attachment resource

#

## Example Usage

```hcl
resource "ksyun_kec_network_interface_attachment" "default" {
  network_interface_id = "ebd74f60-04f1-4b67-91e0-xxxxxxxxxxxx"
  instance_id          = "110d1ce0-113e-4019-8b39-xxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the instance.
* `network_interface_id` - (Required, ForceNew) The ID of the network interface.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `network_interface_type` - The type of the network interface.


## Import

KEC network interface attachment can be imported using the id, e.g.

```
$ terraform import ksyun_kec_network_interface_attachment.default ${network_interface_id}:${instance_id}
```

