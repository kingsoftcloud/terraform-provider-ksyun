variable "available_zone" {
  default = "cn-beijing-6a"
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "foo" {
  subnet_name  = "ksyun-subnet-terraform"
  cidr_block   = "10.7.0.0/21"
  subnet_type  = "Reserve"
  dhcp_ip_from = "10.7.0.2"
  dhcp_ip_to   = "10.7.0.253"
  vpc_id       = "${ksyun_vpc.default.id}"
  gateway_ip   = "10.7.0.1"
  dns1         = "198.18.254.41"
  dns2         = "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_krds_security_group" "krds_sec_group_237" {
  output_file = "output_file"
  security_group_name = "terraform_security_group_237"
  security_group_description = "terraform-security-group-237"
  security_group_rule{
    security_group_rule_protocol = "182.133.0.0/16"
    security_group_rule_name = "asdf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.134.0.0/16"
    security_group_rule_name = "asdf2"
  }
}


resource "ksyun_krds" "houbin_terraform_4"{
  db_instance_class= "db.ram.1|db.disk.20"
  db_instance_name = "houbin_terraform_1-1"
  db_instance_type = "HRDS"
  engine = "mysql"
  engine_version = "5.6"
  master_user_name = "admin"
  master_user_password = "123qweASD123"
  vpc_id = "${ksyun_vpc.default.id}"
  subnet_id = "${ksyun_subnet.foo.id}"
  bill_type = "DAY"
  security_group_id = "${ksyun_krds_security_group.krds_sec_group_237.id}"
  preferred_backup_time = "02:00-03:00"
  port=3306

  # 是否分配EIP（默认为false）
  # instance_has_eip = true

  parameters {
    name = "auto_increment_increment"
    value = "10"
  }

  parameters {
    name = "binlog_format"
    value = "ROW"
  }

  # 部分parameters的更新，需要重启才能生效
  # 将force_restart设置为true，会在修改参数后自动执行重启
  # force_restart = true

}



