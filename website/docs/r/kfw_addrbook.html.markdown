---
subcategory: "KFW"
layout: "ksyun"
page_title: "ksyun: ksyun_kfw_addrbook"
sidebar_current: "docs-ksyun-resource-kfw_addrbook"
description: |-
  Provides a Cloud Firewall Address Book resource.
---

# ksyun_kfw_addrbook

Provides a Cloud Firewall Address Book resource.

#

## Example Usage

```hcl
resource "ksyun_kfw_addrbook" "default" {
  cfw_instance_id = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  addrbook_name   = "test-addrbook"
  ip_version      = "IPv4"
  ip_address      = ["10.1.1.11", "10.2.2.21"]
  description     = "test address book"
}
```

## Argument Reference

The following arguments are supported:

* `addrbook_name` - (Required) Address book name.
* `cfw_instance_id` - (Required, ForceNew) Cloud Firewall Instance ID.
* `ip_address` - (Required) IP addresses.
* `ip_version` - (Required) IP version. Valid values: IPv4, IPv6.
* `description` - (Optional) Description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `addrbook_id` - Address Book ID.
* `citation_count` - Number of references.
* `create_time` - Creation time.


## Import

Cloud Firewall Address Book can be imported using the `addrbook_id`, e.g.

```
$ terraform import ksyun_kfw_addrbook.default addrbook_id
```

