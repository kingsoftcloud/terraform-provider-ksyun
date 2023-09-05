---
subcategory: "ALB"
layout: "ksyun"
page_title: "ksyun: ksyun_alb_listener_cert_groups"
sidebar_current: "docs-ksyun-datasource-alb_listener_cert_groups"
description: |-
  This data source provides a list of ALB listener cert group resources according to their ID.
---

# ksyun_alb_listener_cert_groups

This data source provides a list of ALB listener cert group resources according to their ID.

#

## Example Usage

```hcl
data "ksyun_alb_listener_cert_groups" "default" {
  output_file     = "output_result"
  ids             = []
  alb_listener_id = []
}
```

## Argument Reference

The following arguments are supported:

* `alb_listener_id` - (Optional) One or more ALB Listener IDs.
* `ids` - (Optional) A list of ALB Listener cert group IDs, all the ALB Listener cert groups belong to this region will be retrieved if the ID is `""`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `listener_cert_groups` - An information list of ALB Listener cert groups. Each element contains the following attributes:
  * `alb_listener_cert_group_id` - The ID of the ALB listener cert group.
  * `alb_listener_cert_set` - An information list of ALB Listener certs. Each element contains the following attributes:
    * `cert_authority` - certificate authority.
    * `certificate_id` - The ID of the certificate.
    * `certificate_name` - The name of the certificate.
    * `common_name` - The common name on the certificate.
    * `create_time` - The creation time.
    * `expire_time` - The expire time of the certificate.
  * `alb_listener_id` - The ID of the ALB listener.
  * `id` - ID of the ALB Listener cert group.
* `total_count` - Total number of ALB listener cert groups that satisfy the condition.


