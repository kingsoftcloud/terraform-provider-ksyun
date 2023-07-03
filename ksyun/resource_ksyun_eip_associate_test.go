package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"strings"
	"testing"
)

func TestAccKsyunEipAssociationAssociation_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_eip_associate.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEipAssociationAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipAssociationAssociationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipAssociationAssociationExists("ksyun_eip_associate.foo", &val),
					testAccCheckEipAssociationAssociationAttributes(&val),
				),
			},
		},
	})
}
func testAccCheckEipAssociationAssociationExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("EipAssociationAssociation id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		eipAssociationAssociation := make(map[string]interface{})
		eipAssociationAssociation["AllocationId.1"] = strings.Split(rs.Primary.ID, ":")[0]
		eipAssociationAssociation["ProjectId.1"] = rs.Primary.Attributes["project_id"]
		ptr, err := client.eipconn.DescribeAddresses(&eipAssociationAssociation)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["AddressesSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}
func testAccCheckEipAssociationAssociationAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["AddressesSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("EipAssociationAssociation id is empty")
			}
		}
		return nil
	}
}
func testAccCheckEipAssociationAssociationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_eip_associate" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		eipAssociationAssociation := make(map[string]interface{})
		eipAssociationAssociation["AllocationId"] = strings.Split(rs.Primary.ID, ":")[0]
		ptr, err := client.eipconn.DescribeAddresses(&eipAssociationAssociation)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["AddressesSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				if address, ok := l[0].(map[string]interface{}); ok {
					if address["InstanceId"] == strings.Split(rs.Primary.ID, ":")[1] {
						return fmt.Errorf("EipAssociationAssociation still exist")
					}
				}

			}
		}
	}
	return nil
}

const testAccEipAssociationAssociationConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "vpc-test-tf"
  cidr_block = "10.1.0.0/21"
}
resource "ksyun_lb" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  load_balancer_name = "lb-test-tf"
  type = "public"
  subnet_id = ""
  load_balancer_state = "stop"
  private_ip_address = ""
}
data "ksyun_lines" "default" {
  output_file="output_result1"
  line_name="BGP"
}

# Create an eip
resource "ksyun_eip" "default" {
  line_id ="${data.ksyun_lines.default.lines.0.line_id}"
  band_width =1
  charge_type = "PostPaidByDay"
  purchase_time =1
  project_id=100013
}
resource "ksyun_eip_associate" "foo" {
  allocation_id="${ksyun_eip.default.id}"
  instance_type="Slb"
  instance_id="${ksyun_lb.default.id}"
  network_interface_id=""
}
`
