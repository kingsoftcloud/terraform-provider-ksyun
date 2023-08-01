---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_data_guard_group"
sidebar_current: "docs-ksyun-resource-data_guard_group"
description: |-
  Provides a data guard group resource.
---

# ksyun_data_guard_group

Provides a data guard group resource.

#

## Example Usage

```hcl
resource "ksyun_data_guard_group" "foo" {
  data_guard_name = "your data guard name"
  data_guard_type = "host"
}
```

## Argument Reference

The following arguments are supported:

* `data_guard_name` - (Required) The name of data guard group.
* `data_guard_type` - (Optional) The data guard group display type, Value options: [ host, domain ]. <br> - Host: represent machine level <br> - Domain: represent the domain of disaster tolerance <br> Note: some account doesn't support domain data guard.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `data_guard_capacity` - The capacity of data guard group.
* `data_guard_id` - The id of data guard group.
* `data_guard_level` - The data guard group level, if the value is Host represent machine level, and the tol represent the domain of disaster tolerance.
* `data_guard_used_size` - This data guard group includes the amount of instances.


## Import

Data guard group can be imported using the `id`, e.g.

```
$ terraform import ksyun_data_guard_group.foo "data_guard_id"
```

