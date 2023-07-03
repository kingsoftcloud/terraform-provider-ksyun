---
subcategory: "SQLServer"
layout: "ksyun"
page_title: "ksyun: ksyun_sqlservers"
sidebar_current: "docs-ksyun-datasource-sqlservers"
description: |-
  Query HRDS-ss instance information
---

# ksyun_sqlservers

Query HRDS-ss instance information

#

## Example Usage

```hcl
data "ksyun_sqlservers" "search-sqlservers" {
  output_file            = "output_file"
  db_instance_identifier = "***"
  db_instance_type       = "HRDS-SS"
  keyword                = ""
  order                  = ""
  project_id             = ""
  marker                 = ""
  max_records            = ""
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_identifier` - (Optional) source instance identifier.
* `db_instance_type` - (Optional) HRDS hrds (highly available), RR (read-only), trds (temporary).
* `keyword` - (Optional) fuzzy filter by name / VIP.
* `marker` - (Optional) record start offset.
* `max_records` - (Optional) the maximum number of entries in the result of each page. Value range: 1-100.
* `order` - (Optional) case sensitive, value range: default (default sorting method), group (sorting by replication group, will rank read-only instances after their primary instances).
* `output_file` - (Optional) will return the file name of the content store.
* `project_id` - (Optional) defaults to all projects.
* `sqlservers` - (Optional) a list of instance.

The `db_instance_class` object supports the following:

* `disk` - (Optional) hard disk size..
* `id` - (Optional) id of the DBInstanceClass.
* `iops` - (Optional) IOPS.
* `max_conn` - (Optional) max connection.
* `mem` - (Optional) memory size.
* `ram` - (Optional) memory size.
* `vcpus` - (Optional) the number of the vcpu.

The `db_source` object supports the following:

* `db_instance_identifier` - (Optional) DB instance Identifier.
* `db_instance_name` - (Optional) DB instance name.
* `db_instance_type` - (Optional) DB instance Type.
* `point_in_time` - (Optional) Point in time.

The `read_replica_db_instance_identifiers` object supports the following:

* `id` - (Optional) ID.
* `read_replica_db_instance_identifier` - (Optional) ID.
* `vip` - (Optional) VIP.

The `sqlservers` object supports the following:

* `audit` - (Optional) Audit.
* `availability_zone` - (Optional) AZ.
* `bill_type_id` - (Optional) Bill Type ID.
* `bill_type` - (Optional) Bill type.
* `datastore_version_id` - (Optional) Data store version ID.
* `db_instance_class` - (Optional) instance specification.
* `db_instance_identifier` - (Optional) instance ID.
* `db_instance_name` - (Optional) instance name.
* `db_instance_status` - (Optional) instance status.
* `db_instance_type` - (Optional) instance type.
* `db_parameter_group_id` - (Optional) parameter group ID.
* `db_source` - (Optional) DB source.
* `disk_used` - (Optional) hard disk usage.
* `eip_port` - (Optional) EIP Port number.
* `eip` - (Optional) EIP address.
* `engine_version` - (Optional) database engine version.
* `engine` - (Optional) Database Engine.
* `group_id` - (Optional) group ID.
* `instance_create_time` - (Optional) instance creation time.
* `master_availability_zone` - (Optional) Master AZ.
* `master_user_name` - (Optional) primary account user name.
* `multi_availability_zone` - (Optional) Multi availability zone.
* `order_id` - (Optional) Order ID.
* `order_source` - (Optional) Order source.
* `order_type` - (Optional) Order type.
* `order_use` - (Optional) Order Use.
* `port` - (Optional) Port number.
* `preferred_backup_time` - (Optional) preferred backup time.
* `product_id` - (Optional) Product ID.
* `product_what` - (Optional) Product what.
* `project_id` - (Optional) Project ID.
* `project_name` - (Optional) Project name.
* `publicly_accessible` - (Optional) publicly accessible.
* `read_replica_db_instance_identifiers` - (Optional) read only instance.
* `region` - (Optional) Region.
* `rip` - (Optional) rip.
* `security_group_id` - (Optional) Security group ID.
* `service_end_time` - (Optional) Service end time.
* `service_start_time` - (Optional) Service start time.
* `slave_availability_zone` - (Optional) Slave AZ.
* `sub_order_id` - (Optional) Sub order ID.
* `subnet_id` - (Optional) Subnet ID.
* `vip` - (Optional) virtual IP.
* `vpc_id` - (Optional) virtual private network ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total_count` - Total number of instance that satisfy the condition.


