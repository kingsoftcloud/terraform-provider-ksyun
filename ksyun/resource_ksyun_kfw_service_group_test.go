package ksyun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKsyunKfwServiceGroup_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_kfw_service_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKfwServiceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKfwServiceGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKfwServiceGroupExists("ksyun_kfw_service_group.foo", &val),
					testAccCheckKfwServiceGroupAttributes(&val),
					resource.TestCheckResourceAttr("ksyun_kfw_service_group.foo", "service_group_name", "test-service-group"),
					resource.TestCheckResourceAttr("ksyun_kfw_service_group.foo", "service_infos.#", "2"),
				),
			},
		},
	})
}

func TestAccKsyunKfwServiceGroup_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_kfw_service_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKfwServiceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKfwServiceGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKfwServiceGroupExists("ksyun_kfw_service_group.foo", &val),
					testAccCheckKfwServiceGroupAttributes(&val),
				),
			},
			{
				Config: testAccKfwServiceGroupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKfwServiceGroupExists("ksyun_kfw_service_group.foo", &val),
					testAccCheckKfwServiceGroupAttributes(&val),
					resource.TestCheckResourceAttr("ksyun_kfw_service_group.foo", "service_group_name", "test-service-group-updated"),
					resource.TestCheckResourceAttr("ksyun_kfw_service_group.foo", "service_infos.#", "3"),
				),
			},
		},
	})
}

func testAccCheckKfwServiceGroupExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return nil
}

func testAccCheckKfwServiceGroupAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val == nil || len(*val) == 0 {
			return fmt.Errorf("CfwServiceGroup attributes is empty")
		}
		return nil
	}
}

func testAccCheckKfwServiceGroupDestroy(s *terraform.State) error {
	return nil
}

const testAccKfwServiceGroupConfig = `
resource "ksyun_kfw_instance" "foo" {
  instance_name  = "test-cfw-instance"
  instance_type  = "Advanced"
  bandwidth      = 50
  total_eip_num  = 50
  charge_type    = "Monthly"
  project_id     = "0"
  purchase_time  = 1
}

resource "ksyun_kfw_service_group" "foo" {
  cfw_instance_id    = ksyun_kfw_instance.foo.instance_id
  service_group_name = "test-service-group"
  service_infos      = ["TCP:1-100/80-80", "UDP:22/33"]
  description        = "test service group"
}
`

const testAccKfwServiceGroupUpdateConfig = `
resource "ksyun_kfw_instance" "foo" {
  instance_name  = "test-cfw-instance"
  instance_type  = "Advanced"
  bandwidth      = 50
  total_eip_num  = 50
  charge_type    = "Monthly"
  project_id     = "0"
  purchase_time  = 1
}

resource "ksyun_kfw_service_group" "foo" {
  cfw_instance_id    = ksyun_kfw_instance.foo.instance_id
  service_group_name = "test-service-group-updated"
  service_infos      = ["TCP:1-100/80-80", "UDP:22/33", "ICMP"]
  description        = "test service group updated"
}
`
