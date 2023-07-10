package ksyun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKsyunSubnet_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_subnet.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSubnetDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists("ksyun_subnet.foo", &val),
					testAccCheckSubnetAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunSubnet_update(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_subnet.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSubnetDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists("ksyun_subnet.foo", &val),
					testAccCheckSubnetAttributes(&val),
				),
			},
			{
				Config: testAccSubnetUpdateConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists("ksyun_subnet.foo", &val),
					testAccCheckSubnetAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckSubnetExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("subnet id is empty")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		subnet := make(map[string]interface{})
		subnet["SubnetId.1"] = rs.Primary.ID
		ptr, err := client.vpcconn.DescribeSubnets(&subnet)

		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["SubnetSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}

func testAccCheckSubnetAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["SubnetSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("subnet id is empty")
			}
		}
		return nil
	}
}

func testAccCheckSubnetDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_subnet" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		subnet := make(map[string]interface{})
		subnet["SubnetId.1"] = rs.Primary.ID
		ptr, err := client.vpcconn.DescribeSubnets(&subnet)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["SubnetSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("Subnet still exist")
			}
		}
	}

	return nil
}

const testAccSubnetConfig = `
provider "ksyun" {
	region = "cn-guangzhou-1"
}

data "ksyun_availability_zones" "default" {
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
  provided_ipv6_cidr_block = true
}
resource "ksyun_subnet" "foo" {
  subnet_name      = "ksyun-subnet-tf"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Normal"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  provided_ipv6_cidr_block = true
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}
`

const testAccSubnetUpdateConfig = `
provider "ksyun" {
	region = "cn-guangzhou-1"
}
data "ksyun_availability_zones" "default" {
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
  provided_ipv6_cidr_block = true
}
resource "ksyun_subnet" "foo" {
  subnet_name      = "ksyun-subnet-tf-update"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Normal"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  provided_ipv6_cidr_block = true
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}
`
