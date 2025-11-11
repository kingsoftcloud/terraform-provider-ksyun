---
subcategory: "KFW"
layout: "ksyun"
page_title: "ksyun: ksyun_kfw_addrbooks"
sidebar_current: "docs-ksyun-datasource-kfw_addrbooks"
description: |-
  This data source provides a list of Cloud Firewall Address Book resources according to their instance ID, address book ID, name, and other filters.
---

# ksyun_kfw_addrbooks

This data source provides a list of Cloud Firewall Address Book resources according to their instance ID, address book ID, name, and other filters.

#

## Example Usage

```hcl
data "ksyun_kfw_addrbooks" "default" {
  output_file     = "output_result"
  cfw_instance_id = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  ids             = []
}
```

## Argument Reference

The following arguments are supported:

* `cfw_instance_id` - (Required) Cloud Firewall Instance ID.
* `ids` - (Optional) A list of Address Book IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `kfw_addrbooks` - It is a nested type which documented below.
  * `addrbook_id` - Address Book ID.
  * `addrbook_name` - Address book name.
  * `cfw_instance_id` - Cloud Firewall Instance ID.
  * `citation_count` - Number of references.
  * `create_time` - Creation time.
  * `description` - Description.
  * `ip_address` - IP addresses.
  * `ip_version` - IP version. Valid values: IPv4, IPv6.
* `total_count` - Total number of Address Books that satisfy the condition.


