---
subcategory: "KRDS"
layout: "ksyun"
page_title: "ksyun: ksyun_krds"
sidebar_current: "docs-ksyun-resource-krds"
description: |-
  Provides an RDS instance resource. A DB instance is an isolated database environment in the cloud. A DB instance can contain multiple user-created databases.
---

# ksyun_krds

Provides an RDS instance resource. A DB instance is an isolated database environment in the cloud. A DB instance can contain multiple user-created databases.

#

## Example Usage

```hcl
# Create a RDS MySQL instance

provider "ksyun" {
  region     = "cn-shanghai-3"
  access_key = ""
  secret_key = ""
}

variable "available_zone" {
  default = "cn-shanghai-3a"
}

resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_subnet" "foo" {
  subnet_name       = "ksyun-subnet-tf"
  cidr_block        = "10.7.0.0/21"
  subnet_type       = "Reserve"
  dhcp_ip_from      = "10.7.0.2"
  dhcp_ip_to        = "10.7.0.253"
  vpc_id            = "${ksyun_vpc.default.id}"
  gateway_ip        = "10.7.0.1"
  dns1              = "198.18.254.41"
  dns2              = "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_krds_security_group" "krds_sec_group_14" {
  security_group_name        = "terraform_security_group_14"
  security_group_description = "terraform-security-group-14"
  security_group_rule {
    security_group_rule_protocol = "182.133.0.0/16"
    security_group_rule_name     = "asdf"
  }
  security_group_rule {
    security_group_rule_protocol = "182.134.0.0/16"
    security_group_rule_name     = "asdf2"
  }
}

resource "ksyun_krds" "my_rds_xx" {
  db_instance_class     = "db.ram.2|db.disk.21"
  db_instance_name      = "houbin_terraform_1-n"
  db_instance_type      = "HRDS"
  engine                = "mysql"
  engine_version        = "5.7"
  master_user_name      = "admin"
  master_user_password  = "123qweASD123"
  vpc_id                = "${ksyun_vpc.default.id}"
  subnet_id             = "${ksyun_subnet.foo.id}"
  bill_type             = "DAY"
  security_group_id     = "${ksyun_krds_security_group.krds_sec_group_14.id}"
  preferred_backup_time = "01:00-02:00"
  availability_zone_1   = "cn-shanghai-3a"
  availability_zone_2   = "cn-shanghai-3b"
  port                  = 3306
}

# Create a RDS MySQL instance with specific parameters

provider "ksyun" {
  region     = "cn-shanghai-3"
  access_key = ""
  secret_key = ""
}

variable "available_zone" {
  default = "cn-shanghai-3a"
}

resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_subnet" "foo" {
  subnet_name       = "ksyun-subnet-tf"
  cidr_block        = "10.7.0.0/21"
  subnet_type       = "Reserve"
  dhcp_ip_from      = "10.7.0.2"
  dhcp_ip_to        = "10.7.0.253"
  vpc_id            = "${ksyun_vpc.default.id}"
  gateway_ip        = "10.7.0.1"
  dns1              = "198.18.254.41"
  dns2              = "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_krds_security_group" "krds_sec_group_14" {
  output_file                = "output_file"
  security_group_name        = "terraform_security_group_14"
  security_group_description = "terraform-security-group-14"
  security_group_rule {
    security_group_rule_protocol = "182.133.0.0/16"
    security_group_rule_name     = "asdf"
  }
  security_group_rule {
    security_group_rule_protocol = "182.134.0.0/16"
    security_group_rule_name     = "asdf2"
  }
}

resource "ksyun_krds" "my_rds_xx" {
  output_file           = "output_file"
  db_instance_class     = "db.ram.2|db.disk.21"
  db_instance_name      = "houbin_terraform_1-n"
  db_instance_type      = "HRDS"
  engine                = "mysql"
  engine_version        = "5.7"
  master_user_name      = "admin"
  master_user_password  = "123qweASD123"
  vpc_id                = "${ksyun_vpc.default.id}"
  subnet_id             = "${ksyun_subnet.foo.id}"
  bill_type             = "DAY"
  security_group_id     = "${ksyun_krds_security_group.krds_sec_group_14.id}"
  preferred_backup_time = "01:00-02:00"
  parameters {
    name  = "auto_increment_increment"
    value = "8"
  }

  parameters {
    name  = "binlog_format"
    value = "ROW"
  }

  parameters {
    name  = "delayed_insert_limit"
    value = "108"
  }
  parameters {
    name  = "auto_increment_offset"
    value = "2"
  }
  availability_zone_1 = "cn-shanghai-3a"
  availability_zone_2 = "cn-shanghai-3b"
  instance_has_eip    = true
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_class` - (Required) this value regex db.ram.d{1,9}|db.disk.d{1,9}, db.ram is rds random access memory size, db.disk is disk size.
* `db_instance_name` - (Required) instance name.
* `db_instance_type` - (Required) instance type, valid values: HRDS, TRDS, ERDS, SINGLERDS.
* `engine_version` - (Required) db engine version only support 5.5|5.6|5.7|8.0.
* `engine` - (Required, ForceNew) engine is db type, only support mysql|percona.
* `master_user_name` - (Required, ForceNew) database primary account name.
* `master_user_password` - (Required) master account password.
* `subnet_id` - (Required, ForceNew) ID of the subnet.
* `vpc_id` - (Required, ForceNew) ID of th VPC.
* `availability_zone_1` - (Optional) zone 1.
* `availability_zone_2` - (Optional) zone 2.
* `bill_type` - (Optional, ForceNew) bill type, valid values: DAY, YEAR_MONTH, HourlyInstantSettlement. Default is DAY.
* `db_parameter_template_id` - (Optional, ForceNew) the id of parameter template that for the krds being created.
* `duration` - (Optional) purchase duration in months.
* `force_restart` - (Optional) Set it to true to make some parameter efficient when modifying them. Default to false.
* `instance_has_eip` - (Optional) attach eip for instance.
* `parameters` - (Optional) database parameters.
* `port` - (Optional) port number.
* `preferred_backup_time` - (Optional) backup time.
* `project_id` - (Optional) project ID.
* `security_group_id` - (Optional) proprietary security group id for krds.
* `tags` - (Optional) the tags of the resource.
* `vip` - (Optional) virtual IP.

The `parameters` object supports the following:

* `name` - (Required) name of the parameter.
* `value` - (Required) value of the parameter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `db_instance_identifier` - instance ID.
* `db_parameter_group_id` - ID of the parameter group.
* `eip_port` - EIP port.
* `eip` - EIP address.
* `instance_create_time` - instance create time.
* `region` - region code.


## Import

KRDS can be imported using the id, e.g.

```
$ terraform import ksyun_krds.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

