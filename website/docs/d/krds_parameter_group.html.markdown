---
subcategory: "KRDS"
layout: "ksyun"
page_title: "ksyun: ksyun_krds_parameter_group"
sidebar_current: "docs-ksyun-datasource-krds_parameter_group"
description: |-
  Query ksyun krds parameter group information
---

# ksyun_krds_parameter_group

Query ksyun krds parameter group information

#

## Example Usage

```hcl
provider "ksyun" {
  region = "cn-beijing-6"
}

data "ksyun_krds_parameter_group" "foo" {
  output_file = "output_result"
  // if you give db_parameter_group_id will return the single krds parameter group
  // if you don't give this value, it will return a list of krds parameter groups
  db_parameter_group_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

  // keyword is a filter value that can query the results by name of description
  keyword = "name or description"
}
```

## Argument Reference

The following arguments are supported:

* `db_parameter_group_id` - (Optional) The id of db parameter group.
* `keyword` - (Optional) The keyword uses to filter parameter group.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `db_parameter_groups` - An information list of krds db parameter groups. Each element contains the following attributes:
  * `db_parameter_group_id` - The krds db parameter group id.
  * `db_parameter_group_name` - The krds db parameter group name.
  * `description` - The description of this db parameter group.
  * `engine_version` - The version of engine.
  * `engine` - The db parameter group adapts to what krds engine.
  * `parameters` - The custom parameters.
* `total_count` - Total number of snapshot policies resources that satisfy the condition.


