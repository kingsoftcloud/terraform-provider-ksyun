---
subcategory: "PostgreSQL"
layout: "ksyun"
page_title: "ksyun: ksyun_postgresql_rr"
sidebar_current: "docs-ksyun-resource-postgresql_rr"
description: |-
  Provides a PostgreSQL Read Only instance resource. A DB read only instance is an isolated database environment in the cloud.
---

# ksyun_postgresql_rr

Provides a PostgreSQL Read Only instance resource. A DB read only instance is an isolated database environment in the cloud.

#

## Example Usage

```hcl
resource "ksyun_postgresql_rr" "my_postgresql_rr" {
  db_instance_identifier = "******"
  db_instance_class      = "db.ram.2|db.disk.50"
  db_instance_name       = "houbin_terraform_888_rr_1"
  bill_type              = "DAY"
  security_group_id      = "******"

  parameters {
    name  = "auto_increment_increment"
    value = "7"
  }

  parameters {
    name  = "binlog_format"
    value = "ROW"
  }
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_class` - (Required) this value regex db.ram.d{1,3}|db.disk.d{1,5}, db.ram is PostgreSQL random access memory size, db.disk is disk size.
* `db_instance_identifier` - (Required, ForceNew) passes in the instance ID of the PostgreSQL highly available instance. A PostgreSQL highly available instance can have at most three read-only instances.
* `db_instance_name` - (Required) instance name.
* `availability_zone_1` - (Optional, ForceNew) zone 1.
* `bill_type` - (Optional, ForceNew) bill type, valid values: DAY, YEAR_MONTH, HourlyInstantSettlement. Default is DAY.
* `duration` - (Optional) purchase duration in months.
* `force_restart` - (Optional) Set it to true to make some parameter efficient when modifying them. Default to false.
* `instance_has_eip` - (Optional) attach eip for instance.
* `parameters` - (Optional) database parameters.
* `port` - (Optional) port number.
* `project_id` - (Optional) project ID.
* `security_group_id` - (Optional) proprietary security group id for postgresql.
* `tags` - (Optional) the tags of the resource.
* `vcpus` - (Optional, ForceNew) The number of vCPUs for the DB instance. If not specified, defaults to half of the memory size.
* `vip` - (Optional) virtual IP.

The `parameters` object supports the following:

* `name` - (Required) name of the parameter.
* `value` - (Required) value of the parameter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `db_instance_type` - instance type, valid values: HRDS_PG (highly available), RR_PG (read-only).
* `db_parameter_group_id` - ID of the parameter group.
* `eip_port` - EIP port.
* `eip` - EIP address.
* `engine_version` - db engine version, valid values: 10|11|12.5|13|14|15|17.
* `engine` - engine is db type, postgresql.
* `instance_create_time` - instance create time.
* `region` - region code.


## Import

PostgreSQL Read Only instance resource can be imported using the id, e.g.

```
$ terraform import ksyun_postgresql_rr.my_postgresql_rr 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

