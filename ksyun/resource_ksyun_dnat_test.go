package ksyun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKsyunDnat_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_dnat.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnatConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnatExists("ksyun_dnat.foo", &val),
					testAccCheckDnatAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckDnatExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Dnat id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		Dnat := make(map[string]interface{})
		Dnat["DnatId.1"] = rs.Primary.ID
		ptr, err := client.vpcconn.DescribeDnats(&Dnat)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["DnatSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}
		*val = *ptr
		return nil
	}
}
func testAccCheckDnatAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["DnatSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("Dnat id is empty")
			}
		}
		return nil
	}
}
func testAccCheckDnatDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_dnat" {
			client := testAccProvider.Meta().(*KsyunClient)
			Dnat := make(map[string]interface{})
			Dnat["DnatId.1"] = rs.Primary.ID
			ptr, err := client.vpcconn.DescribeDnats(&Dnat)

			// Verify the error is what we want
			if err != nil {
				return err
			}
			if ptr != nil {
				l := (*ptr)["DnatSet"].([]interface{})
				if len(l) == 0 {
					continue
				} else {
					return fmt.Errorf("Dnat still exist")
				}
			}
		}

	}

	return nil
}

const testAccDnatConfig = `
provider "ksyun" {
	region = "cn-guangzhou-1"
}

data "ksyun_images" "centos-7_5" {
  platform= "centos-7.5"
}
data "ksyun_availability_zones" "default" {
}

resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.foo.id}"
  security_group_name="ksyun-security-group-nat"
}

resource "ksyun_instance" "foo" {
  image_id="${data.ksyun_images.centos-7_5.images.0.image_id}"
  instance_type="N3.2B"

  #max_count=1
  #min_count=1
  subnet_id="${ksyun_subnet.foo.id}"
  instance_password="Xuan663222"
  keep_image_login=false
  charge_type="Daily"
  purchase_time=1
  security_group_id=["${ksyun_security_group.default.id}"]
  instance_name="ksyun-kec-tf-nat"
  sriov_net_support="false"
  project_id=100012
}

resource "ksyun_nat" "foo" {
  nat_name = "ksyun-nat-tf"
  nat_mode = "Subnet"
  nat_type = "public"
  band_width = 1
  charge_type = "DailyPaidByTransfer"
  vpc_id = "${ksyun_vpc.foo.id}"
  nat_ip_number = 3 
}
resource "ksyun_vpc" "foo" {
	vpc_name        = "tf-vpc-nat"
	cidr_block = "10.0.5.0/24"
}

resource "ksyun_subnet" "foo" {
  subnet_name      = "tf-acc-nat-subnet1"
  cidr_block = "10.0.5.0/24"
  subnet_type = "Normal"
  vpc_id  = "${ksyun_vpc.foo.id}"
  gateway_ip = "10.0.5.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}

resource "ksyun_dnat" "foo" {
	nat_id = ksyun_nat.foo.id
	dnat_name = "tf_dnat"
	ip_protocol = "Any"
	nat_ip =ksyun_nat.foo.nat_ip_set[0].nat_ip
	private_ip_address = ksyun_instance.foo.private_ip_address
	description = "test"
	private_port = "Any"
	public_port = "Any"
}
`
