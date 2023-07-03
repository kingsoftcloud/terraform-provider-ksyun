---
subcategory: "SSH key"
layout: "ksyun"
page_title: "ksyun: ksyun_ssh_key"
sidebar_current: "docs-ksyun-resource-ssh_key"
description: |-
  Provides an SSH key resource.
---

# ksyun_ssh_key

Provides an SSH key resource.

#

## Example Usage

```hcl
resource "ksyun_ssh_key" "default" {
  key_name   = "ssh_key_tf"
  public_key = "ssh-rsa xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `key_name` - (Optional) name of the key.
* `public_key` - (Optional, ForceNew) public key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - creation time of the key.
* `key_id` - ID of the key.
* `private_key` - private key.


## Import

SSH key can be imported using the id, e.g.

```
$ terraform import ksyun_ssh_key.default xxxxxxxxxxxx
```

