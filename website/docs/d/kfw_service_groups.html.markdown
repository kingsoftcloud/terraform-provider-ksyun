---
subcategory: "KFW"
layout: "ksyun"
page_title: "ksyun: ksyun_kfw_service_groups"
sidebar_current: "docs-ksyun-datasource-kfw_service_groups"
description: |-
  This data source provides a list of Cloud Firewall Service Group resources according to their instance ID, service group ID, name, and other filters.
---

# ksyun_kfw_service_groups

This data source provides a list of Cloud Firewall Service Group resources according to their instance ID, service group ID, name, and other filters.

#

## Example Usage

```hcl
data "ksyun_kfw_service_groups" "default" {
  output_file     = "output_result"
  cfw_instance_id = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  ids             = []
}
```

## Argument Reference

The following arguments are supported:

* `cfw_instance_id` - (Required) Cloud Firewall Instance ID.
* `ids` - (Optional) A list of Service Group IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `kfw_service_groups` - It is a nested type which documented below.
  * `cfw_instance_id` - Cloud Firewall Instance ID.
  * `citation_count` - Number of references.
  * `description` - Description.
  * `service_group_id` - Service Group ID.
  * `service_group_name` - Service group name.
  * `service_infos` - Service information. Format: protocol:src_port_min-src_port_max/dest_port_min-dest_port_max. Example: TCP:1-100/70-80, UDP:22/33, ICMP.
* `total_count` - Total number of Service Groups that satisfy the condition.


