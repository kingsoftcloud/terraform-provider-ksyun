---
subcategory: "PostgreSQL"
layout: "ksyun"
page_title: "ksyun: ksyun_postgresql_parameter_group"
sidebar_current: "docs-ksyun-resource-postgresql_parameter_group"
description: |-
  Provides a PostgreSQL parameter template group.
---

# ksyun_postgresql_parameter_group

Provides a PostgreSQL parameter template group.

#

## Example Usage

```hcl
resource "ksyun_postgresql_parameter_group" "dpg" {
  name           = "tf_dpg_on_hcl"
  description    = "tf_configuration_test"
  engine         = "postgresql"
  engine_version = "15"
  parameters {
    name  = "auto_increment_increment"
    value = "8"
  }
  parameters {
    name  = "binlog_format"
    value = "ROW"
  }
  parameters {
    name  = "delayed_insert_limit"
    value = "108"
  }
  parameters {
    name  = "auto_increment_offset"
    value = "2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, ForceNew) postgresql database version, valid values: 10|11|12.5|13|14|15|17.
* `engine` - (Required, ForceNew) postgresql database engine.
* `name` - (Required) the name of postgresql parameter group.
* `description` - (Optional) The description of this db parameter group.
* `parameters` - (Optional) database parameters.

The `parameters` object supports the following:

* `name` - (Required) name of the parameter.
* `value` - (Required) value of the parameter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `db_parameter_group_id` - The id of postgresql parameter group.
* `resource_name` - identify this resource.


## Import

Tag can be imported using the `id`, e.g.

```
$ terraform import ksyun_postgresql_parameter_group.foo "id"
```

