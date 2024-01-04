---
subcategory: "KCR"
layout: "ksyun"
page_title: "ksyun: ksyun_kcrs_instances"
sidebar_current: "docs-ksyun-datasource-kcrs_instances"
description: |-
  This data source provides a list of instance resources according to their id.
---

# ksyun_kcrs_instances

This data source provides a list of instance resources according to their id.

## Example Usage

```hcl
data "ksyun_kcrs_instances" "foo" {
  output_file = "kcrs_instance_output_result"
  ids         = []
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Kcrs Instance IDs, all the Kcrs Instances belong to this region will be retrieved if the ID is `""`.
* `name_regex` - (Optional) A regex string to filter results by instance name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `project_ids` - (Optional) One or more project IDs. If its value is none, returns instance all of project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - It is a nested type which documented below.
  * `charge_type` - Charge Type.
  * `create_time` - Created Time.
  * `instance_id` - ID of the instance.
  * `instance_name` - name of the instance.
  * `instance_status` - Instance status.
  * `instance_type` - type of the instance.
  * `internal_endpoint` - Internal endpoint dns.
  * `project_id` - Project Id.
* `total_count` - Total number of kcrs instances that satisfy the condition.


