---
subcategory: "KFW"
layout: "ksyun"
page_title: "ksyun: ksyun_kfw_service_group"
sidebar_current: "docs-ksyun-resource-kfw_service_group"
description: |-
  Provides a Cloud Firewall Service Group resource.
---

# ksyun_kfw_service_group

Provides a Cloud Firewall Service Group resource.

#

## Example Usage

```hcl
resource "ksyun_kfw_service_group" "default" {
  cfw_instance_id    = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  service_group_name = "test-service-group"
  service_infos      = ["TCP:1-65535/1-65535", "UDP:22/33"]
  description        = "test service group"
}
```

## Argument Reference

The following arguments are supported:

* `cfw_instance_id` - (Required, ForceNew) Cloud Firewall Instance ID.
* `service_group_name` - (Required) Service group name.
* `service_infos` - (Required) Service information. Format: protocol:src_port_min-src_port_max/dest_port_min-dest_port_max. Example: TCP:1-100/80-80, UDP:22/33, ICMP.
* `description` - (Optional) Description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `citation_count` - Number of references.
* `service_group_id` - Service Group ID.


## Import

Cloud Firewall Service Group can be imported using the `service_group_id`, e.g.

```
$ terraform import ksyun_kfw_service_group.default service_group_id
```

