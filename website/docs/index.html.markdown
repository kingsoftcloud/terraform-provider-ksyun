---
layout: "ksyun"
page_title: "Provider: Ksyun"
sidebar_current: "docs-ksyun-index"
description: |-
  The Ksyun provider is used to interact with many resources supported by Ksyun. The provider needs to be configured with the proper credentials before it can be used.
---

# Ksyun Provider

~> **NOTE:** This guide requires an available Ksyun account or sub-account with project to create resources.

The Ksyun provider is used to interact with the
resources supported by Ksyun. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Ksyun Provider
provider "ksyun" {
   access_key = "your ak"
   secret_key = "your sk"
   region = "cn-beijing-6"
}
data "ksyun_availability_zones" "default" {

}
# Query image
data "ksyun_images" "default" {
  output_file="output_result"
  ids=[]
  name_regex="centos-7.0-20180927115228"
  is_public=true
  image_source="system"
}
#Create vpc
resource "ksyun_vpc" "default" {
  vpc_name   = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}
#Create subnet
resource "ksyun_subnet" "default" {
  subnet_name = "${var.subnet_name}"
  cidr_block = "10.1.0.0/21"
  subnet_type = "Normal"
  dhcp_ip_from = "10.1.0.2"
  dhcp_ip_to = "10.1.0.253"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.1.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}
#Create security group
resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="${var.security_group_name}"
}
#Create ssh key
resource "ksyun_ssh_key" "default" {
  key_name="ssh_key_tf"
  public_key=""
}
# Create instance 
resource "ksyun_instance" "default" {
  image_id="${data.ksyun_images.centos-7_5.images.0.image_id}"
  instance_type="S4.1A"
  system_disk{
    disk_type="SSD3.0"
    disk_size=30
  }
  data_disk_gb=0
  data_disk =[
    {
      type="SSD3.0"
      size=20
      delete_with_instance=true
    }
  ]
  subnet_id="${ksyun_subnet.default.id}"
  instance_password="Xuan663222"
  keep_image_login=false
  charge_type="Daily"
  purchase_time=1
  security_group_id=["${ksyun_security_group.default.id}","${ksyun_security_group.default2.id}"]
  private_ip_address=""
  instance_name="xuan-tf-combine"
  instance_name_suffix=""
  sriov_net_support=false
  project_id=0
  data_guard_id=""
  key_id=["${ksyun_ssh_key.default.id}"]
  force_delete=true
}


```

## Authentication

The Ksyun provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Static credentials
- Environment variables

### Static credentials

Static credentials can be provided by adding an `public_key` and `private_key` in-line in the
Ksyun provider block:

Usage:

```hcl
provider "ksyun" {
   access_key = "your ak"
   secret_key = "your sk"
   region = "cn-beijing-6"
}
```

### Environment variables

You can provide your credentials via `KSYUN_ACCESS_KEY` and `KSYUN_SECRET_KEY`
environment variables, representing your Ksyun public key and private key respectively.
`KSYUN_REGION` is also used, if applicable:

```hcl
provider "ksyun" {}
```

Usage:

```shell
$ export KSYUN_ACCESS_KEY="your_public_key"
$ export KSYUN_SECRET_KEY="your_private_key"
$ export KSYUN_REGION="cn-beijing-6"

$ terraform plan
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the Ksyun
 `provider` block:

* `access_key` - (Required) This is the Ksyun public key. It must be provided, but
  it can also be sourced from the `KSYUN_ACCESS_KEY` environment variable.

* `secret_key` - (Required) This is the Ksyun private key. It must be provided, but
  it can also be sourced from the `KSYUN_SECRET_KEY` environment variable.

* `region` - (Required) This is the Ksyun region. It must be provided, but
  it can also be sourced from the `KSYUN_REGION` environment variables.

* `max_retries` - (Optional) This is the max retry attempts number. Default max retry attempts number is `3`.

* `insecure` - (Optional) This is a switch to disable/enable https. (Default: `false`, means enable https).

* `domain` - (Optional) This is the base url of KSYUN API endpoint. (Default: `api.ksyun.com`) Setup to corresponding base URL if you are using private cloud or other delicated regions. 

* `endpoint` - (Optional) This is the access domain name for KS3 services.

* `dry_run` - (Optional, Boolean) Whether enable `dry_run` while operating SDK. 

* `ignore_service` - (Optional, Boolean) Whether ignore customer's service. 

* `force_https` - (Optional, Boolean) Force use https protocol for communication between sdk and remote server.

* `http_keepalive` - (Optional, Boolean) Whether use http keepalive, if false, disables HTTP keep-alives and will only use the connection to the server for a single HTTP request. 

* `http_proxy` - (Optional) Indicating a http proxy server that the cyber traffic via. 

## Testing

Credentials must be provided via the `KSYUN_ACCESS_KEY`, `KSYUN_SECRET_KEY` environment variables in order to run acceptance tests.
