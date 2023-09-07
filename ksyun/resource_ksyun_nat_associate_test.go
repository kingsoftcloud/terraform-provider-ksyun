package ksyun

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func TestAccKsyunNatAssociation_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_nat_associate.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNatAssociationDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccNatAssociationConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatAssociationExists("ksyun_nat_associate.foo", &val),
					testAccCheckNatAssociationAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckNatAssociationExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf(" Nat id is empty ")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		Nat := make(map[string]interface{})
		Nat["NatId.1"] = strings.Split(rs.Primary.ID, ":")[0]
		projectErr := getProjectInfo(&Nat, client)
		if projectErr != nil {
			return projectErr
		}
		ptr, err := client.vpcconn.DescribeNats(&Nat)

		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["NatSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}

func testAccCheckNatAssociationAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["NatSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf(" Nat id is empty ")
			}
		}
		return nil
	}
}

func testAccCheckNatAssociationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_nat_associate" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		Nat := make(map[string]interface{})
		Nat["NatId.1"] = strings.Split(rs.Primary.ID, ":")[0]
		subnetId := strings.Split(rs.Primary.ID, ":")[1]
		projectErr := getProjectInfo(&Nat, client)
		if projectErr != nil {
			return projectErr
		}
		ptr, err := client.vpcconn.DescribeNats(&Nat)
		logger.Debug(logger.ReqFormat, "DescribeNats", ptr)
		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["NatSet"].([]interface{})
			if len(l) == 1 {
				flag := true
				if nat, ok := l[0].(map[string]interface{}); ok {
					if nat["AssociateNatSet"] == nil {
						continue
					}
					if associates, o1 := nat["AssociateNatSet"].([]interface{}); o1 {
						for _, v := range associates {
							if subnet, o2 := v.(map[string]interface{}); o2 {
								if subnetId == subnet["SubnetId"].(string) {
									flag = false
									break
								}
							}
						}
					}
					if flag {
						continue
					} else {
						return fmt.Errorf(" Nat Associate Still Exist ")
					}
				}
			} else {
				continue
			}
		}
	}

	return nil
}

const testAccNatAssociationConfig = `

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

resource "ksyun_nat_associate" "foo" {
  nat_id = "${ksyun_nat.foo.id}"
  network_interface_id = "${ksyun_instance.foo.network_interface_id}"
}
`
