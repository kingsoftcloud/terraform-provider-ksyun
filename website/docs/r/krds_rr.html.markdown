---
subcategory: "KRDS"
layout: "ksyun"
page_title: "ksyun: ksyun_krds_rr"
sidebar_current: "docs-ksyun-resource-krds_rr"
description: |-
  Provides an RDS Read Only instance resource. A DB read only instance is an isolated database environment in the cloud.
---

# ksyun_krds_rr

Provides an RDS Read Only instance resource. A DB read only instance is an isolated database environment in the cloud.

#

## Example Usage

```hcl
resource "ksyun_krds_rr" "my_rds_rr" {
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

* `db_instance_class` - (Required) this value regex db.ram.d{1,3}|db.disk.d{1,5}, db.ram is rds random access memory size, db.disk is disk size.
* `db_instance_identifier` - (Required, ForceNew) passes in the instance ID of the RDS highly available instance. A RDS highly available instance can have at most three read-only instances.
* `db_instance_name` - (Required) instance name.
* `availability_zone_1` - (Optional, ForceNew) zone 1.
* `bill_type` - (Optional, ForceNew) bill type, valid values: DAY, YEAR_MONTH, HourlyInstantSettlement. Default is DAY.
* `duration` - (Optional) purchase duration in months.
* `force_restart` - (Optional) Set it to true to make some parameter efficient when modifying them. Default to false.
* `instance_has_eip` - (Optional) attach eip for instance.
* `parameters` - (Optional) database parameters.
* `port` - (Optional) port number.
* `project_id` - (Optional) project ID.
* `security_group_id` - (Optional) proprietary security group id for krds.
* `tags` - (Optional) the tags of the resource.
* `vip` - (Optional) virtual IP.

The `parameters` object supports the following:

* `name` - (Required) name of the parameter.
* `value` - (Required) value of the parameter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `db_instance_type` - instance type, valid values: HRDS, TRDS, ERDS, SINGLERDS.
* `db_parameter_group_id` - ID of the parameter group.
* `eip_port` - EIP port.
* `eip` - EIP address.
* `engine_version` - db engine version only support 5.5|5.6|5.7|8.0.
* `engine` - engine is db type, only support mysql|percona.
* `instance_create_time` - instance create time.
* `region` - region code.


## Import

RDS Read Only instance resource can be imported using the id, e.g.

```
$ terraform import ksyun_krds_rr.my_rds_rr 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

