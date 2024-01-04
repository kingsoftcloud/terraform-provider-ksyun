---
subcategory: "KCR"
layout: "ksyun"
page_title: "ksyun: ksyun_kcrs_instance"
sidebar_current: "docs-ksyun-resource-kcrs_instance"
description: |-
  Provides a Kcrs Repository Instance resource.
---

# ksyun_kcrs_instance

Provides a Kcrs Repository Instance resource.

## Example Usage

```hcl
# Create a Kcrs Repository Instance
resource "ksyun_kcrs_instance" "foo" {
  instance_name = "tfunittest"
  instance_type = "basic"
}

# Create a Kcrs Repository Instance and open public access
resource "ksyun_kcrs_instance" "foo" {
  instance_name         = "tfunittest"
  instance_type         = "basic"
  open_public_operation = true

  # open public access with external policy that permits an address, ip or cidr, to access this repository
  external_policy {
    entry = "192.168.2.133"
    desc  = "ddd"
  }
  external_policy {
    entry = "192.168.2.123/32"
    desc  = "ddd"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, ForceNew) Repository instance name.
* `instance_type` - (Required, ForceNew) The type of instance. Valid Values: 'basic', 'premium'.
* `charge_type` - (Optional, ForceNew) Charge type of the instance. Valid Values: 'HourlyInstantSettlement'. Default: 'HourlyInstantSettlement'.
* `delete_bucket` - (Optional) Whether delete bucket with this instance is removing.
* `external_policy` - (Optional) The external access policy. It's activated when 'open_public_operation' is true.
* `open_public_operation` - (Optional) Control public network access.
* `project_id` - (Optional) The id of the project.

The `external_policy` object supports the following:

* `entry` - (Required) External policy entry. Submit to CIDR or IP.
* `desc` - (Optional) The external policy description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `expired_time` - Expired time.
* `instance_status` - Repository instance status.
* `internal_endpoint` - Internal endpoint address.
* `public_domain` - Public domain.


## Import

KcrsInstance can be imported using the id, e.g.

```
$ terraform import ksyun_kcrs_instance.foo 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

