package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunEip_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_eip.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("ksyun_eip.foo", &val),
					testAccCheckEipAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunEip_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_eip.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("ksyun_eip.foo", &val),
					testAccCheckEipAttributes(&val),
				),
			},
			{
				Config: testAccEipUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("ksyun_eip.foo", &val),
					testAccCheckEipAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckEipExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Eip id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		eip := make(map[string]interface{})
		eip["AllocationId.1"] = rs.Primary.ID
		eip["ProjectId.1"] = rs.Primary.Attributes["project_id"]
		ptr, err := client.eipconn.DescribeAddresses(&eip)
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
func testAccCheckEipAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["AddressesSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("Eip id is empty")
			}
		}
		return nil
	}
}
func testAccCheckEipDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_eip" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		eip := make(map[string]interface{})
		eip["AllocationId.1"] = rs.Primary.ID
		ptr, err := client.eipconn.DescribeAddresses(&eip)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["AddressesSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("Eip still exist")
			}
		}
	}

	return nil
}

const testAccEipConfig = `
data "ksyun_lines" "default" {
  output_file="output_result1"
  line_name="BGP"
}
resource "ksyun_eip" "foo" {
  line_id =data.ksyun_lines.default.lines.0.line_id
  band_width =1
  charge_type = "PostPaidByDay"
  purchase_time =1
  project_id=100012
}
`

const testAccEipUpdateConfig = `
data "ksyun_lines" "default" {
  output_file="output_result1"
  line_name="BGP"
}
# Create an eip
resource "ksyun_eip" "foo" {
  line_id ="${data.ksyun_lines.default.lines.0.line_id}"
  band_width =10
  charge_type = "PostPaidByDay"
  purchase_time =1
  project_id=0
}
`
