---
subcategory: "RabbitMQ"
layout: "ksyun"
page_title: "ksyun: ksyun_rabbitmq_security_rule"
sidebar_current: "docs-ksyun-resource-rabbitmq_security_rule"
description: |-
  Provides a Rabbitmq Security Rule resource.
---

# ksyun_rabbitmq_security_rule

Provides a Rabbitmq Security Rule resource.

#

## Example Usage

```hcl
resource "ksyun_rabbitmq_security_rule" "default" {
  instance_id = "InstanceId"
  cidr        = "192.168.10.1/32"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The id of the rabbitmq instance.
* `cidr` - (Optional, ForceNew) network cidr.
* `cidrs` - (Optional, **Deprecated**) `cidrs` is deprecated use resourceKsyunRabbitmq.cidrs instead  network cidrs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



