---
subcategory: "KRDS"
layout: "ksyun"
page_title: "ksyun: ksyun_krds"
sidebar_current: "docs-ksyun-datasource-krds"
description: |-
  Query HRDS and RDS-rr instance information
---

# ksyun_krds

Query HRDS and RDS-rr instance information

#

## Example Usage

```hcl
data "ksyun_krds" "search-krds" {
  output_file            = "output_file"
  db_instance_identifier = "***"
  db_instance_type       = "HRDS,RR,TRDS"
  keyword                = ""
  order                  = ""
  project_id             = ""
  marker                 = ""
  max_records            = ""
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_identifier` - (Optional) instance ID (passed in the instance ID to get the details of the instance, otherwise get the list).
* `db_instance_status` - (Optional) status of the instance, ACTIVE or INVALID.
* `db_instance_type` - (Optional) HRDS (highly available), RR (read-only), TRDS (temporary).
* `keyword` - (Optional) fuzzy filter by name / VIP.
* `marker` - (Optional) record start offset.
* `max_records` - (Optional) the maximum number of entries in the result of each page. Value range: 1-100.
* `order` - (Optional) case sensitive, value range: default (default sorting method), group (sorting by replication group, will rank read-only instances after their primary instances).
* `output_file` - (Optional) will return the file name of the content store.
* `project_id` - (Optional) the default value is all projects.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `krds` - An information list of KRDS.
  * `availability_zone` - Availability zone.
  * `db_instance_class` - instance specification.
    * `disk` - hard disk size.
    * `id` - id of the DBInstanceClass.
    * `iops` - IOPS.
    * `max_conn` - max connection.
    * `mem` - memory size.
    * `ram` - memory size.
    * `vcpus` - number of CPUs.
  * `db_instance_identifier` - instance ID.
  * `db_instance_status` - instance status.
  * `disk_used` - hard disk usage.
  * `eip_port` - EIP port.
  * `eip` - EIP address.
  * `group_id` - ID of the parameter group.
  * `instance_create_time` - instance creation time.
  * `product_id` - Product ID.
  * `region` - Region.
  * `service_end_time` - service end time.
  * `service_start_time` - Service start time.
  * `vip` - virtual IP.
* `total_count` - Total number of resources that satisfy the condition.


