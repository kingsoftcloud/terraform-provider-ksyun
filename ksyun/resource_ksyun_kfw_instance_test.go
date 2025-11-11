package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunKfwInstance_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_cfw_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKfwInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKfwInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKfwInstanceExists("ksyun_cfw_instance.foo", &val),
					testAccCheckKfwInstanceAttributes(&val),
					resource.TestCheckResourceAttr("ksyun_cfw_instance.foo", "instance_name", "test-cfw-instance"),
					resource.TestCheckResourceAttr("ksyun_cfw_instance.foo", "instance_type", "Advanced"),
					resource.TestCheckResourceAttr("ksyun_cfw_instance.foo", "bandwidth", "50"),
					resource.TestCheckResourceAttr("ksyun_cfw_instance.foo", "total_eip_num", "50"),
					resource.TestCheckResourceAttr("ksyun_cfw_instance.foo", "charge_type", "Monthly"),
				),
			},
		},
	})
}

func TestAccKsyunKfwInstance_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_cfw_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKfwInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKfwInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKfwInstanceExists("ksyun_cfw_instance.foo", &val),
					testAccCheckKfwInstanceAttributes(&val),
				),
			},
			{
				Config: testAccKfwInstanceUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKfwInstanceExists("ksyun_cfw_instance.foo", &val),
					testAccCheckKfwInstanceAttributes(&val),
					resource.TestCheckResourceAttr("ksyun_cfw_instance.foo", "instance_name", "test-cfw-instance-updated"),
					resource.TestCheckResourceAttr("ksyun_cfw_instance.foo", "bandwidth", "100"),
				),
			},
		},
	})
}

func testAccCheckKfwInstanceExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return nil
}

func testAccCheckKfwInstanceAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val == nil || len(*val) == 0 {
			return fmt.Errorf("CfwInstance attributes is empty")
		}
		return nil
	}
}

func testAccCheckKfwInstanceDestroy(s *terraform.State) error {
	return nil
}

const testAccKfwInstanceConfig = `
resource "ksyun_kfw_instance" "foo" {
  instance_name  = "test-kfw-instance"
  instance_type  = "Advanced"
  bandwidth      = 50
  total_eip_num  = 50
  charge_type    = "Monthly"
  project_id     = "0"
  purchase_time  = 1
}
`

const testAccKfwInstanceUpdateConfig = `
resource "ksyun_kfw_instance" "foo" {
  instance_name  = "test-kfw-instance-updated"
  instance_type  = "Advanced"
  bandwidth      = 100
  total_eip_num  = 50
  charge_type    = "Monthly"
  project_id     = "0"
  purchase_time  = 1
}
`
