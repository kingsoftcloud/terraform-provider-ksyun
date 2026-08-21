---
subcategory: "PostgreSQL"
layout: "ksyun"
page_title: "ksyun: ksyun_postgresql"
sidebar_current: "docs-ksyun-resource-postgresql"
description: |-
  Provides a PostgreSQL instance resource. A DB instance is an isolated database environment in the cloud. A DB instance can contain multiple user-created databases.
---

# ksyun_postgresql

Provides a PostgreSQL instance resource. A DB instance is an isolated database environment in the cloud. A DB instance can contain multiple user-created databases.

#

## Example Usage

```hcl
# Create a PostgreSQL instance

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

resource "ksyun_postgresql_security_group" "postgresql_sec_group_14" {
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

resource "ksyun_postgresql" "my_postgresql_xx" {
  db_instance_class     = "db.ram.2|db.disk.21"
  db_instance_name      = "houbin_terraform_1-n"
  db_instance_type      = "HRDS_PG"
  engine                = "postgresql"
  engine_version        = "15"
  master_user_name      = "admin"
  master_user_password  = "123qweASD123"
  vpc_id                = "${ksyun_vpc.default.id}"
  subnet_id             = "${ksyun_subnet.foo.id}"
  bill_type             = "DAY"
  security_group_id     = "${ksyun_postgresql_security_group.postgresql_sec_group_14.id}"
  preferred_backup_time = "01:00-02:00"
  availability_zone_1   = "cn-shanghai-3a"
  availability_zone_2   = "cn-shanghai-3b"
  port                  = 5432
}

# Create a PostgreSQL instance with specific parameters

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

resource "ksyun_postgresql_security_group" "postgresql_sec_group_14" {
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

resource "ksyun_postgresql" "my_postgresql_xx" {
  output_file           = "output_file"
  db_instance_class     = "db.ram.2|db.disk.21"
  db_instance_name      = "houbin_terraform_1-n"
  db_instance_type      = "HRDS_PG"
  engine                = "postgresql"
  engine_version        = "15"
  master_user_name      = "admin"
  master_user_password  = "123qweASD123"
  vpc_id                = "${ksyun_vpc.default.id}"
  subnet_id             = "${ksyun_subnet.foo.id}"
  bill_type             = "DAY"
  security_group_id     = "${ksyun_postgresql_security_group.postgresql_sec_group_14.id}"
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

* `db_instance_class` - (Required) this value regex db.ram.d{1,9}|db.disk.d{1,9}, db.ram is postgresql random access memory size, db.disk is disk size.
* `db_instance_name` - (Required) instance name.
* `db_instance_type` - (Required, ForceNew) instance type, valid values: HRDS_PG (highly available), RR_PG (read-only).
* `engine_version` - (Required) db engine version, valid values: 10|11|12.5|13|14|15|17.
* `engine` - (Required, ForceNew) engine is db type, postgresql.
* `master_user_name` - (Required, ForceNew) database primary account name.
* `master_user_password` - (Required) master account password.
* `subnet_id` - (Required, ForceNew) ID of the subnet.
* `vpc_id` - (Required, ForceNew) ID of th VPC.
* `availability_zone_1` - (Optional) zone 1.
* `availability_zone_2` - (Optional) zone 2.
* `bill_type` - (Optional, ForceNew) bill type, valid values: DAY, YEAR_MONTH, HourlyInstantSettlement. Default is DAY.
* `db_parameter_template_id` - (Optional, ForceNew) the id of parameter template that for the postgresql being created.
* `duration` - (Optional) purchase duration in months.
* `force_restart` - (Optional) Set it to true to make some parameter efficient when modifying them. Default to false.
* `instance_has_eip` - (Optional) attach eip for instance.
* `parameters` - (Optional) database parameters.
* `port` - (Optional) port number.
* `preferred_backup_time` - (Optional) backup time.
* `project_id` - (Optional) project ID.
* `security_group_id` - (Optional) proprietary security group id for postgresql.
* `tags` - (Optional) the tags of the resource.
* `vcpus` - (Optional, ForceNew) The number of vCPUs for the DB instance. If not specified, defaults to half of the memory size.
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

PostgreSQL can be imported using the id, e.g.

```
$ terraform import ksyun_postgresql.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

