---
subcategory: "KCR"
layout: "ksyun"
page_title: "ksyun: ksyun_kcrs_vpc_attachment"
sidebar_current: "docs-ksyun-resource-kcrs_vpc_attachment"
description: |-
  Provides an internal access attachment resource with a vpc under kcrs repository instance.
---

# ksyun_kcrs_vpc_attachment

Provides an internal access attachment resource with a vpc under kcrs repository instance.

## Example Usage

```hcl
# repository instance
resource "ksyun_kcrs_instance" "foo" {
  instance_name = "tfunittest"
  instance_type = "basic"
}

# To attach a vpc for an instance
resource "ksyun_kcrs_vpc_attachment" "foo" {
  instance_id           = ksyun_kcrs_instance.foo.id
  vpc_id                = "vpc_id"
  reserve_subnet_id     = "subnet_id"
  enable_vpc_domain_dns = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) Instance id of repository.
* `reserve_subnet_id` - (Required, ForceNew) The id of subnet type is '**Reserve**'.
* `vpc_id` - (Required, ForceNew) Vpc id.
* `enable_vpc_domain_dns` - (Optional) Whether to enable vpc domain dns. Default value is `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `dns_parse_status` - Status of the DNS parsed.
* `eni_lb_ip` - IP address of the internal access.
* `internal_endpoint_dns` - Endpoint Domain of the internal access.
* `status` - Status of the internal access.


## Import

KcrsVpcAttachment can be imported using `instance_id:vpc_id`, e.g.

```
$ terraform import ksyun_kcrs_vpc_attachment.foo ${instance_id}:${vpc_id}
```

