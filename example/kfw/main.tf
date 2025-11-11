terraform {
  required_providers {
    ksyun = {
      source = "kingsoftcloud/ksyun"
      version = ">= 1.16.0"
    }
  }
}


# Create an ksyun_kfw
resource "ksyun_kfw_instance" "foo" {
  instance_name  = "test-kfw-instance"
  instance_type  = "Advanced"
  bandwidth      = 50
  total_eip_num  = 50
  charge_type    = "Monthly"
  project_id     = "0"
  purchase_time  = 1
}

resource "ksyun_kfw_acl" "foo" {
  cfw_instance_id   = ksyun_kfw_instance.foo.id
  acl_name         = "test-acl-rule"
  direction        = "in"
  src_type         = "ip"
  src_ips          = ["1.1.1.1", "2.2.2.2"]
  dest_type        = "ip"
  dest_ips         = ["3.3.3.3"]
  service_type     = "service"
  service_infos     = ["TCP:1-100/80-80"]
  app_type         = "any"
  policy           = "accept"
  status           = "start"
  priority_position = "after+1"
  description      = "test acl rule"
}

resource "ksyun_kfw_addrbook" "foo" {
  cfw_instance_id   = ksyun_kfw_instance.foo.id
  addrbook_name   = "test-addrbook"
  ip_version      = "IPv4"
  ip_address      = ["1.1.1.1", "2.2.2.2"]
  description     = "test address book"
}

resource "ksyun_kfw_service_group" "foo" {
  cfw_instance_id    = ksyun_kfw_instance.foo.id
  service_group_name = "test-service-group"
  service_infos      = ["TCP:1-100/80-80", "UDP:22/33"]
  description        = "test service group"
}

data "ksyun_kfw_acls" "default" {
  output_file       = "output_result"
  cfw_instance_id   = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  ids               = []
}

data "ksyun_kfw_addrbooks" "default" {
  output_file       = "output_result"
  cfw_instance_id   = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  ids               = []
}

data "ksyun_kfw_instances" "default" {
  output_file = "output_result"
  ids         = []
}

data "ksyun_kfw_service_groups" "default" {
  output_file     = "output_result"
  cfw_instance_id = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  ids             = []
}

