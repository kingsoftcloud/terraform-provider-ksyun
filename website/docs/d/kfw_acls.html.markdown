---
subcategory: "KFW"
layout: "ksyun"
page_title: "ksyun: ksyun_kfw_acls"
sidebar_current: "docs-ksyun-datasource-kfw_acls"
description: |-
  This data source provides a list of Cloud Firewall ACL Rule resources according to their instance ID, ACL ID, name, and other filters.
---

# ksyun_kfw_acls

This data source provides a list of Cloud Firewall ACL Rule resources according to their instance ID, ACL ID, name, and other filters.

#

## Example Usage

```hcl
data "ksyun_kfw_acls" "default" {
  output_file     = "output_result"
  cfw_instance_id = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  ids             = []
}
```

## Argument Reference

The following arguments are supported:

* `cfw_instance_id` - (Required) Cloud Firewall Instance ID.
* `ids` - (Optional) A list of ACL Rule IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `kfw_acls` - It is a nested type which documented below.
  * `acl_id` - ACL Rule ID.
  * `acl_name` - ACL rule name.
  * `app_type` - Application type. Valid values: app, any.
  * `app_value` - Application values.
  * `cfw_instance_id` - Cloud Firewall Instance ID.
  * `create_time` - Creation time.
  * `description` - Description.
  * `dest_addrbooks` - Destination address book IDs.
  * `dest_host` - Destination domain names.
  * `dest_hostbook` - Destination host book IDs.
  * `dest_ips` - Destination IP addresses.
  * `dest_type` - Destination address type. Valid values: ip, addrbook, any.
  * `direction` - Direction. Valid values: in (inbound), out (outbound).
  * `hit_count` - Hit count.
  * `policy` - Action. Valid values: accept, deny.
  * `priority_position` - Priority position. Format: after+priority or before+priority. Example: after+1, before+1.
  * `priority` - Priority value.
  * `service_groups` - Service group IDs.
  * `service_infos` - Service information. Format: protocol:src_port_min-src_port_max/dest_port_min-dest_port_max. Example: TCP:1-100/80-80, UDP:22/33, ICMP.
  * `service_type` - Service type. Valid values: service, servicegroup, any.
  * `src_addrbooks` - Source address book IDs.
  * `src_ips` - Source IP addresses.
  * `src_type` - Source address type. Valid values: ip, addrbook, zone, any.
  * `src_zones` - Source zones (geographic regions).
  * `status` - Status. Valid values: start, stop.
* `total_count` - Total number of ACL Rules that satisfy the condition.


