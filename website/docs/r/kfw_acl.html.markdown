---
subcategory: "KFW"
layout: "ksyun"
page_title: "ksyun: ksyun_kfw_acl"
sidebar_current: "docs-ksyun-resource-kfw_acl"
description: |-
  Provides a Cloud Firewall ACL Rule resource.
---

# ksyun_kfw_acl

Provides a Cloud Firewall ACL Rule resource.

#

## Example Usage

```hcl
resource "ksyun_kfw_acl" "default" {
  cfw_instance_id   = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  acl_name          = "test-acl-rule"
  direction         = "in"
  src_type          = "ip"
  src_ips           = ["10.0.0.11", "10.0.0.21"]
  dest_type         = "ip"
  dest_ips          = ["10.0.0.31"]
  service_type      = "service"
  service_infos     = ["TCP:1-65535/1-65535"]
  app_type          = "any"
  policy            = "accept"
  status            = "start"
  priority_position = "after1"
  description       = "test acl rule"
}
```

## Argument Reference

The following arguments are supported:

* `acl_name` - (Required) ACL rule name.
* `app_type` - (Required) Application type. Valid values: app, any.
* `cfw_instance_id` - (Required, ForceNew) Cloud Firewall Instance ID.
* `dest_type` - (Required) Destination address type. Valid values: ip, addrbook, any.
* `direction` - (Required, ForceNew) Direction. Valid values: in (inbound), out (outbound).
* `policy` - (Required) Action. Valid values: accept, deny.
* `priority_position` - (Required) Priority position. Format: after+priority or before+priority. Example: after+1, before+1.
* `service_type` - (Required) Service type. Valid values: service, servicegroup, any.
* `src_type` - (Required) Source address type. Valid values: ip, addrbook, zone, any.
* `status` - (Required) Status. Valid values: start, stop.
* `app_value` - (Optional) Application values.
* `description` - (Optional) Description.
* `dest_addrbooks` - (Optional) Destination address book IDs.
* `dest_host` - (Optional) Destination domain names.
* `dest_hostbook` - (Optional) Destination host book IDs.
* `dest_ips` - (Optional) Destination IP addresses.
* `service_groups` - (Optional) Service group IDs.
* `service_infos` - (Optional) Service information. Format: protocol:src_port_min-src_port_max/dest_port_min-dest_port_max. Example: TCP:1-100/80-80, UDP:22/33, ICMP.
* `src_addrbooks` - (Optional) Source address book IDs.
* `src_ips` - (Optional) Source IP addresses.
* `src_zones` - (Optional) Source zones (geographic regions).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `acl_id` - ACL Rule ID.
* `create_time` - Creation time.
* `hit_count` - Hit count.
* `priority` - Priority value.


## Import

Cloud Firewall ACL Rule can be imported using the `acl_id`, e.g.

```
$ terraform import ksyun_cfw_acl.default acl_id
```

