---
subcategory: "SSH key"
layout: "ksyun"
page_title: "ksyun: ksyun_ssh_keys"
sidebar_current: "docs-ksyun-datasource-ssh_keys"
description: |-
  This data source provides a list of SSH Key resources according to their SSH Key ID.
---

# ksyun_ssh_keys

This data source provides a list of SSH Key resources according to their SSH Key ID.

#

## Example Usage

```hcl
data "ksyun_ssh_keys" "default" {
  output_file = "output_result"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of SSH Key IDs, all the SSH Key belong to this region will be retrieved if the ID is `""`.
* `key_name` - (Optional) ssh key name.
* `key_names` - (Optional) a list of ssh key name.
* `name_regex` - (Optional) A regex string to filter results by kay name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `keys` - An information list of ssh key. Each element contains the following attributes:
  * `create_time` - creation time of the key.
  * `key_id` - ID of the key.
  * `key_name` - name of the key.
  * `public_key` - public key.
* `total_count` - Total number of keys that satisfy the condition.


