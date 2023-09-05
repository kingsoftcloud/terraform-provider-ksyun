---
subcategory: "ALB"
layout: "ksyun"
page_title: "ksyun: ksyun_alb_listener_cert_group"
sidebar_current: "docs-ksyun-resource-alb_listener_cert_group"
description: |-
  Provides a ALB Listener cert group resource.
---

# ksyun_alb_listener_cert_group

Provides a ALB Listener cert group resource.

#

## Example Usage

```hcl
# network and security group configuration
resource "ksyun_vpc" "example" {
  vpc_name   = "tf-alb-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "example" {
  subnet_name       = "tf-alb-subnet"
  cidr_block        = "10.0.1.0/24"
  subnet_type       = "Normal"
  availability_zone = var.az
  vpc_id            = ksyun_vpc.example.id
}

resource "ksyun_security_group" "example" {
  vpc_id              = ksyun_vpc.example.id
  security_group_name = "tf_sg"
}

resource "ksyun_security_group_entry" "example" {
  security_group_id = ksyun_security_group.example.id
  cidr_block        = "0.0.0.0/0"
  protocol          = "ip"
  direction         = "in"
}

# ---------------------------------------------
# resource ksyun alb
resource "ksyun_alb" "example" {
  alb_name    = "tf-alb-example-1"
  alb_version = "standard"
  alb_type    = "public"
  state       = "start"
  charge_type = "PrePaidByHourUsage"
  vpc_id      = ksyun_vpc.example.id
  project_id  = 0
  enabled_log = false
  ip_version  = "ipv4"
}

# query your certificates on ksyun
data "ksyun_certificates" "listener_cert" {
  name_regex = "test"
}

resource "ksyun_alb_listener" "example" {
  alb_id             = ksyun_alb.example.id
  alb_listener_name  = "alb-example-listener"
  protocol           = "HTTPS"
  port               = 8099
  alb_listener_state = "start"
  certificate_id     = data.ksyun_certificates.listener_cert.certificates.0.certificate_id
  http_protocol      = "HTTP1.1"
  session {
    cookie_type                = "ImplantCookie"
    cookie_name                = "KLBRSIDdad"
    session_state              = "start"
    session_persistence_period = 3100
  }
}

resource "ksyun_alb_listener_cert_group" "default" {
  alb_listener_id = ksyun_alb_listener.example.id
  certificate {
    certificate_id = data.ksyun_certificates.listener_cert.certificates.0.certificate_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `alb_listener_id` - (Required, ForceNew) The ID of the ALB Listener.
* `certificate` - (Optional) The certificate included in the cert group.

The `certificate` object supports the following:

* `certificate_id` - (Required) The ID of the certificate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `alb_listener_cert_group_id` - The ID of the ALB Listener Cert Group.


## Import

ALB Listener Cert Group can be imported using the `id`, e.g.

```
$ terraform import ksyun_alb_listener_cert_group.example vserver-abcdefg
```

