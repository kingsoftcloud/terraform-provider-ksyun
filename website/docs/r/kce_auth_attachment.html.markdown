---
subcategory: "KCE"
layout: "ksyun"
page_title: "ksyun: ksyun_kce_auth_attachment"
sidebar_current: "docs-ksyun-resource-kce_auth_attachment"
description: |-
  Provides a kce_auth_attachment resource.
---

# ksyun_kce_auth_attachment

Provides a kce_auth_attachment resource.

## Example Usage

```hcl
resource "ksyun_kce_auth_attachment" "auth" {
  sub_user_id = "38435"
  permissions {
    cluster_id   = "4cf5b24b-de39-4f55-b0ce-fd7b28cb964c"
    cluster_role = "kce:dev"
    namespace    = ""
  }
}
```

## Argument Reference

The following arguments are supported:

* `permissions` - (Required) the permissions of the sub user.
* `sub_user_id` - (Required, ForceNew) the id of the sub user.

The `permissions` object supports the following:

* `cluster_id` - (Required) The id of the kce cluster.
* `cluster_role` - (Required) the role for the sub user in the cluster. Valid Values: kce:admin, kce:dev, kce:ops, kce:restricted, kce:ns:dev, kce:ns:restricted.
* `namespace` - (Optional) the namespace of the cluster role, if it's empty, this authorization will apply in all of namespace.
* `region` - (Optional) the region of the kce cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

kce can be imported using the id, e.g.

```
$ terraform import ksyun_kce_auth_attachment.auth ${sub_user_id}
```

