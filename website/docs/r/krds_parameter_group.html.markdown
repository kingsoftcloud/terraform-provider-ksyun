---
subcategory: "KRDS"
layout: "ksyun"
page_title: "ksyun: ksyun_krds_parameter_group"
sidebar_current: "docs-ksyun-resource-krds_parameter_group"
description: |-
  Provides a krds parameter template groups.
---

# ksyun_krds_parameter_group

Provides a krds parameter template groups.

#

## Example Usage

```hcl
resource "ksyun_krds_parameter_group" "dpg" {
  name           = "tf_dpg_on_hcl"
  description    = "tf_configuration_test"
  engine         = "mysql"
  engine_version = "5.7"
  parameters = {
    connect_timeout         = 20
    max_prepared_stmt_count = 65535
    max_user_connections    = 65535
    back_log                = 65535
    innodb_open_files       = 1000
  }
}
```

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required) krds database version. Value options:<br> - Mysql: [ 5.5, 5.6, 5.7, 8.0 ] <br> - Percona: [ 5.6 ] <br> - Consistent_mysql: [ 5.7 ] <br> - Ebs_mysql: [ 5.6, 5.7 ].
* `engine` - (Required) krds database type. Value options: mysql|percona|consistent_mysql|ebs_mysql.
* `name` - (Required) the name of krds parameter group.
* `description` - (Optional) The description of this db parameter group.
* `parameters` - (Optional) The custom parameters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `db_parameter_group_id` - The id of krds parameter group.
* `resource_name` - identify this resource.


## Import

Tag can be imported using the `id`, e.g.

```
$ terraform import ksyun_krds_parameter_group.foo "id"
```

