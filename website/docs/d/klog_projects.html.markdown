---
subcategory: "KLog"
layout: "ksyun"
page_title: "ksyun: ksyun_klog_projects"
sidebar_current: "docs-ksyun-datasource-klog_projects"
description: |-
  This data source provides a list of KLOG projects.
---

# ksyun_klog_projects

This data source provides a list of KLOG projects.

#

## Example Usage

```hcl
data "ksyun_klog_projects" "default" {
  project_name = "test"
  description  = "online"
  page         = "0"
  size         = "20"
  output_file  = "output_result"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of project.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `page` - (Optional) Page number start from 0.
* `project_name` - (Optional) The name of project.
* `size` - (Optional) Page size, 1 - 500.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `projects` - Project list.
  * `create_time` - When the project was created.
  * `description` - The description of project.
  * `iam_project_id` - The IAMProjectId of project.
  * `iam_project_name` - The IAMProjectName of project.
  * `log_pool_num` - The log pool count of project.
  * `project_name` - The name of project.
  * `region` - The region of project.
  * `status` - The status of project.
  * `tags` - Tags of project.
    * `key` - The key of tag.
    * `value` - The value of tag.
  * `update_time` - When the project was updated.
* `total_count` - Total count of project.


