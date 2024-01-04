---
subcategory: "KCR"
layout: "ksyun"
page_title: "ksyun: ksyun_kcrs_namespace"
sidebar_current: "docs-ksyun-resource-kcrs_namespace"
description: |-
  Provides a namespace resource under kcrs repository instance.
---

# ksyun_kcrs_namespace

Provides a namespace resource under kcrs repository instance.

## Example Usage

```hcl
# repository instance
resource "ksyun_kcrs_instance" "foo" {
  instance_name = "tfunittest"
  instance_type = "basic"
}

# Create a namespace under the repository instance
resource "ksyun_kcrs_namespace" "foo" {
  instance_id = ksyun_kcrs_instance.foo.id
  namespace   = "tftest"
  public      = false
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, ForceNew) The name of namespace.
* `public` - (Required) Whether to be public this namespace.
* `instance_id` - (Optional) Instance id of repository.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

KcrsNamespace can be imported using `instance_id:namespace_name`, e.g.

```
$ terraform import ksyun_kcrs_namespace.foo ${instance_id}:${namespace_name}
```

