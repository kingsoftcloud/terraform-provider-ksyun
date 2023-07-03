---
subcategory: "SLB"
layout: "ksyun"
page_title: "ksyun: ksyun_lb"
sidebar_current: "docs-ksyun-resource-lb"
description: |-
  Provides a Load Balancer resource.
---

# ksyun_lb

Provides a Load Balancer resource.

#

## Example Usage

```hcl
resource "ksyun_lb" "default" {
  vpc_id             = "74d0a45b-472d-49fc-84ad-221e21ee23aa"
  load_balancer_name = "tf-xun1"
  type               = "public"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The ID of the VPC linked to the Load Balancers.
* `access_logs_enabled` - (Optional) Default is `false`, Setting the value to `true` to enable the service.
* `access_logs_s3_bucket` - (Optional) Bucket for storing access logs.
* `ip_version` - (Optional, ForceNew) IP version, valid values: 'all', 'ipv4', 'ipv6'.
* `load_balancer_name` - (Optional) The name of the load balancer.
* `load_balancer_state` - (Optional) The Load Balancers state.Valid Values:'start', 'stop'.
* `private_ip_address` - (Optional, ForceNew) The internal Load Balancers can set an private ip address in Reserve Subnet.
* `project_id` - (Optional) ID of the project.
* `subnet_id` - (Optional, ForceNew) The id of the subnet.only Internal type is Required.
* `tags` - (Optional) the tags of the resource.
* `type` - (Optional, ForceNew) The type of load balancer.Valid Values:'public', 'internal'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time of creation for load balancer.
* `is_waf` - whether it is a waf LB or not.
* `load_balancer_id` - ID of the LB.
* `public_ip` - The IP address of Public IP. It is `""` if `internal` is `true`.
* `state` - associate or disassociate.


## Import

LB can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb.example fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```

