package ksyun

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKsyunKrds_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_krds.rds_terraform_3",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKrdsDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccKrdsConfig,

				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckKrdsExists("ksyun_krds.rds_terraform_3", &val),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccKrdsUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckKrdsExists("ksyun_krds.rds_terraform_3", &val),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckKrdsExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found : %s", n)
		}
		if res.Primary.ID == "" {
			return fmt.Errorf("instance is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		req := map[string]interface{}{
			"DBInstanceIdentifier": res.Primary.ID,
		}
		resp, err := client.krdsconn.DescribeDBInstances(&req)
		if err != nil {
			return err
		}
		if resp != nil {
			bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
			if !dataOk {
				return fmt.Errorf("error on reading Instance(krds)  %+v", (*resp)["Error"])
			}
			instances := bodyData["Instances"].([]interface{})
			if len(instances) == 0 {
				return fmt.Errorf("no instance find, instance number is 0")
			}
		}
		*val = *resp
		return nil
	}
}

func testAccCheckKrdsDestroy(s *terraform.State) error {
	for _, res := range s.RootModule().Resources {
		if res.Type != "ksyun_krds" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		req := map[string]interface{}{
			"DBInstanceIdentifier": res.Primary.ID,
		}
		_, err := client.krdsconn.DescribeDBInstances(&req)
		if err != nil {
			if err.(awserr.Error).Code() == "NOT_FOUND" {
				return nil
			} else if notFoundError(err) {
				return nil
			}
			return err
		}
	}

	return nil
}

const testAccKrdsConfig = `
provider "ksyun" {
	region =  "cn-beijing-6"
}

variable "available_zone" {
  default = "cn-beijing-6a"
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf1"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "foo" {
  subnet_name      = "ksyun-subnet-tf"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Reserve"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_krds_parameter_group" "dpg" {
  name  = "tf_krdpg_on_hcl_with_partial_parameters"
  description    = "acceptance-test"
  engine = "mysql"
  engine_version = "5.5"
  parameters = {
    back_log = 65535
	connect_timeout = 30
  }
}

resource "ksyun_krds" "rds_terraform_3" {
  db_instance_class= "db.ram.1|db.disk.15"
  db_instance_name = "terraform_3"
  db_instance_type = "HRDS"
  engine = "mysql"
  engine_version = "5.5"
  master_user_name = "admin"
  master_user_password = "123qweASD123"
  vpc_id = "${ksyun_vpc.default.id}"
  subnet_id = "${ksyun_subnet.foo.id}"
  bill_type = "DAY"
  preferred_backup_time = "01:00-02:00"
  instance_has_eip = true
  db_parameter_template_id = "${ksyun_krds_parameter_group.dpg.id}"
  parameters {
	name = "innodb_open_files"
	value = 900
  }
  force_restart = true
}

resource "ksyun_tag" "test_tag" {
    key = "exist_tag"
    value = "exist_tag_value3"
    resource_type = "kcs-instance"
    resource_id = "${ksyun_krds.rds_terraform_3.id}"
}
`
const testAccKrdsEipConfig = `
variable "available_zone" {
  default = "cn-beijing-6a"
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf1"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "foo" {
  subnet_name      = "ksyun-subnet-tf"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Reserve"
  dhcp_ip_from = "10.7.0.2"
  dhcp_ip_to = "10.7.0.253"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_krds" "rds_terraform_3"{
  db_instance_class= "db.ram.1|db.disk.15"
  db_instance_name = "terraform_1"
  db_instance_type = "HRDS"
  engine = "mysql"
  engine_version = "5.5"
  master_user_name = "admin"
  master_user_password = "123qweASD123"
  vpc_id = "${ksyun_vpc.default.id}"
  subnet_id = "${ksyun_subnet.foo.id}"
  bill_type = "DAY"
  instance_has_eip = false
  preferred_backup_time = "01:00-02:00"

}

resource "ksyun_tag" "test_tag" {
    key = "exist_tag"
    value = "exist_tag_value4"
    resource_type = "kcs-instance"
    resource_id = "${ksyun_krds.rds_terraform_3.id}"
}
`

const testAccKrdsUpdateConfig = `
provider "ksyun" {
	region =  "cn-beijing-6"
}

variable "available_zone" {
  default = "cn-beijing-6a"
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf1"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "foo" {
  subnet_name      = "ksyun-subnet-tf"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Reserve"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_krds_parameter_group" "dpg" {
  name  = "tf_krdpg_on_hcl_with_partial_parameters"
  description    = "acceptance-test"
  engine = "mysql"
  engine_version = "5.5"
  parameters = {
    back_log = 65535
	connect_timeout = 30
  }
}

resource "ksyun_krds_parameter_group" "dpg_with_apply_2" {
  name  = "tf_krdpg_on_hcl_with_2"
  description    = "acceptance-test-2"
  engine = "mysql"
  engine_version = "5.5"
  parameters = {
    back_log = 65535
	connect_timeout = 30
	long_query_time = 300
	innodb_open_files = 1000
  }
}

resource "ksyun_krds" "rds_terraform_3" {
  db_instance_class= "db.ram.1|db.disk.15"
  db_instance_name = "terraform_3"
  db_instance_type = "HRDS"
  engine = "mysql"
  engine_version = "5.5"
  master_user_name = "admin"
  master_user_password = "123qweASD123"
  vpc_id = "${ksyun_vpc.default.id}"
  subnet_id = "${ksyun_subnet.foo.id}"
  bill_type = "DAY"
  preferred_backup_time = "01:00-02:00"
  instance_has_eip = true
  force_restart = true
  db_parameter_template_id = "${ksyun_krds_parameter_group.dpg.id}"
  parameters  {
  	name = "back_log"
	value = 23344
  }
  parameters  {
  	name = "innodb_open_files"
	value = 700
  }
}

resource "ksyun_tag" "test_tag" {
    key = "exist_tag"
    value = "exist_tag_value3"
    resource_type = "kcs-instance"
    resource_id = "${ksyun_krds.rds_terraform_3.id}"
}
`
