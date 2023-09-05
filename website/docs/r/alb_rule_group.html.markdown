---
subcategory: "ALB"
layout: "ksyun"
page_title: "ksyun: ksyun_alb_rule_group"
sidebar_current: "docs-ksyun-resource-alb_rule_group"
description: |-
  Provides a ALB rule group resource.
---

# ksyun_alb_rule_group

Provides a ALB rule group resource.

#

## Example Usage

```hcl
# network and security group configuration
resource "ksyun_vpc" "test" {
  vpc_name   = "tf-alb-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "test" {
  subnet_name       = "tf-alb-subnet"
  cidr_block        = "10.0.1.0/24"
  subnet_type       = "Normal"
  availability_zone = var.az
  vpc_id            = ksyun_vpc.test.id
}

resource "ksyun_security_group" "test" {
  vpc_id              = ksyun_vpc.test.id
  security_group_name = "tf_sg"
}

resource "ksyun_security_group_entry" "test" {
  security_group_id = ksyun_security_group.test.id
  cidr_block        = "0.0.0.0/0"
  protocol          = "ip"
  direction         = "in"
}

# --------------------------------------------------------
# alb-rule-group relational configuration

# ksyun alb configuration
resource "ksyun_alb" "test" {
  alb_name    = "tf-alb-test"
  alb_version = "standard"
  alb_type    = "public"
  state       = "start"
  charge_type = "PrePaidByHourUsage"
  vpc_id      = ksyun_vpc.test.id
  project_id  = 0
  enabled_log = false
  ip_version  = "ipv4"
}

# ksyun alb listener configuration
resource "ksyun_alb_listener" "test" {
  alb_id             = ksyun_alb.test.id
  alb_listener_name  = "alb-test-listener"
  protocol           = "HTTP"
  port               = 8088
  alb_listener_state = "start"
  http_protocol      = "HTTP1.1"
  session {
    cookie_type                = "ImplantCookie"
    cookie_name                = "KLBRSIDdad"
    session_state              = "start"
    session_persistence_period = 3100
  }
}

# --------------------------------------------
# backend server group and kec instance configuration
# backend server group configuration
resource "ksyun_lb_backend_server_group" "test" {
  backend_server_group_name = "tf_bsg"
  vpc_id                    = ksyun_vpc.test.id
  backend_server_group_type = "Server"
}
resource "ksyun_lb_register_backend_server" "default" {
  backend_server_group_id = ksyun_lb_backend_server_group.test.id
  backend_server_ip       = ksyun_instance.test.0.private_ip_address
  backend_server_port     = 8090
  weight                  = 10
}
resource "ksyun_lb_register_backend_server" "default2" {
  backend_server_group_id = ksyun_lb_backend_server_group.test.id
  backend_server_ip       = ksyun_instance.test.1.private_ip_address
  backend_server_port     = 8090
  weight                  = 10
}

# kec instance creating
data "ksyun_images" "default" {
  output_file  = "output_result"
  name_regex   = "centos-7.0"
  is_public    = true
  image_source = "system"
}

data "ksyun_ssh_keys" "test" {

}

resource "ksyun_instance" "test" {
  count = 2
  security_group_id = [
    ksyun_security_group.test.id
  ]
  subnet_id = ksyun_subnet.test.id
  key_id    = [data.ksyun_ssh_keys.test.keys.0.key_id]

  instance_type = "S6.1A"
  charge_type   = "Daily"
  instance_name = "tf-alb-test-vm"
  project_id    = 0

  image_id = data.ksyun_images.default.images.0.image_id
}

# ksyun_alb_rule_group configuration
resource "ksyun_alb_rule_group" "default" {
  alb_listener_id         = ksyun_alb_listener.test.id
  alb_rule_group_name     = "tf_alb_rule_group"
  backend_server_group_id = ksyun_lb_backend_server_group.test.id
  alb_rule_set {
    alb_rule_type  = "url"
    alb_rule_value = "/test/path"
  }
  alb_rule_set {
    alb_rule_type  = "domain"
    alb_rule_value = "www.ksyun.com"
  }
  listener_sync = "on"
}
```

## Argument Reference

The following arguments are supported:

* `alb_listener_id` - (Required, ForceNew) The ID of the ALB listener.
* `alb_rule_set` - (Required) Rule set, define strategies for being load-balance of backend server.
* `backend_server_group_id` - (Required) The ID of the backend server group.
* `listener_sync` - (Required) Whether to synchronize the health check, session persistence, and load balancing algorithm of the listener. valid values: 'on', 'off'.
* `alb_rule_group_name` - (Optional) The name of the ALB rule group.
* `cookie_name` - (Optional) The name of cookie. Should set it value, when `listener_sync` is off and `cookie_type` is `RewriteCookie`.
* `cookie_type` - (Optional) The type of cookie, valid values: 'ImplantCookie','RewriteCookie'.
* `health_check_state` - (Optional) Status maintained by health examination.Valid Values:'start', 'stop'. Should set it value, when `listener_sync` is off.
* `healthy_threshold` - (Optional) Health threshold.Valid Values:1-10. Should set it value, when `listener_sync` is off.
* `host_name` - (Optional) The service host name of the health check, which is available only for the HTTP or HTTPS health check. Should set it value, when `listener_sync` is off.
* `interval` - (Optional) Interval of health examination.Valid Values:1-3600. Should set it value, when `listener_sync` is off.
* `method` - (Optional) Forwarding mode of listener. Valid Values:'RoundRobin', 'LeastConnections'.
* `session_persistence_period` - (Optional) Session hold timeout. Valid Values:1-86400. Should set it value, when `listener_sync` is off.
* `session_state` - (Optional) The state of session. Valid Values:'start', 'stop'. Should set it value, when `listener_sync` is off.
* `timeout` - (Optional) Health check timeout.Valid Values:1-3600. Should set it value, when `listener_sync` is off.
* `unhealthy_threshold` - (Optional) Unhealthy threshold.Valid Values:1-10. Should set it value, when `listener_sync` is off.
* `url_path` - (Optional) Link to HTTP type listener health check. Should set it value, when `listener_sync` is off.

The `alb_rule_set` object supports the following:

* `alb_rule_type` - (Required) Rule type. valid values: 'domain', 'url'.
* `alb_rule_value` - (Required) Rule value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `alb_rule_group_id` - The ID of the rule group.


## Import

`ksyun_alb_rule_group` can be imported using the id, e.g.

```
$ terraform import ksyun_alb_rule_group.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

