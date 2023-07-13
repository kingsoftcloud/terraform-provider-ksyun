package ksyun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	_ "github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKsyunNetworkInterface_basic(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_kec_network_interface.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkInterfaceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccKniConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_kec_network_interface.foo"),
				),
				ExpectNonEmptyPlan: true,
			},
			// to test terraform when its configuration changes
			{
				Config: testAccKniUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_kec_network_interface.foo"),
					// resource.TestCheckResourceAttr("ksyun_kec_network_interface.foo", "secondary_private_ips.4.ip", "10.7.0.242"),
				),
			},
			// {
			// 	Config: testAccKniEmptyConfig,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckIDExists("ksyun_kec_network_interface.foo"),
			// 	),
			// },
		},
	})
}

func testAccCheckNetworkInterfaceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_kec_network_interface" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		conn := client.vpcconn
		params := make(map[string]interface{})
		params["NetworkInterfaceId.0"] = rs.Primary.ID
		sdkResponse, err := conn.DescribeNetworkInterfaces(&params)

		// Verify the error is what we want
		if err != nil {
			if notFoundError(err) {
				return nil
			}
			return err
		}
		if sdkResponse == nil && len(*sdkResponse) == 0 {
			return fmt.Errorf("sdk response is nil")
		}
		resRaw, _ := getSdkValue("NetworkInterfaceSet", *sdkResponse)
		res, _ := If2Slice(resRaw)
		if len(res) > 0 {
			return fmt.Errorf("kni is still exist")
		}
	}

	return nil
}

const testAccKniConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}

variable "available_zone" {
  default = "cn-beijing-6a"
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf-kni"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "foo" {
  subnet_name      = "ksyun-subnet-tf-kni"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Normal"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="ksyun-security-group-kni"
}

resource "ksyun_kec_network_interface" "foo" {
  network_interface_name   = "tf_kni_secondary_ip"
  subnet_id = ksyun_subnet.foo.id
  security_group_ids = ["${ksyun_security_group.default.id}"]
  // secondary_private_ip_address_count = 5
  secondary_private_ips {
    ip = "10.7.3.136"
  }
  secondary_private_ips {
    ip = "10.7.3.200"
  }
  secondary_private_ips {
    ip = "10.7.0.242"
  }
}
`

const testAccKniUpdateConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}

variable "available_zone" {
  default = "cn-beijing-6a"
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf-kni"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "foo" {
  subnet_name      = "ksyun-subnet-tf-kni"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Normal"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="ksyun-security-group-kni"
}

resource "ksyun_kec_network_interface" "foo" {
  network_interface_name   = "tf_kni_secondary_ip"
  subnet_id = ksyun_subnet.foo.id
  security_group_ids = ["${ksyun_security_group.default.id}"]
  // secondary_private_ip_address_count = 1
  secondary_private_ips {
    ip = "10.7.3.136"
  }
  secondary_private_ips {
    ip = "10.7.3.220"
  }
  secondary_private_ips {
    ip = "10.7.4.232"
  }
}
`
