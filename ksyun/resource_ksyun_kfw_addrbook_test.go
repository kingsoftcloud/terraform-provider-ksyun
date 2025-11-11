package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunKfwAddrbook_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_kfw_addrbook.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKfwAddrbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKfwAddrbookConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKfwAddrbookExists("ksyun_kfw_addrbook.foo", &val),
					testAccCheckKfwAddrbookAttributes(&val),
					resource.TestCheckResourceAttr("ksyun_kfw_addrbook.foo", "addrbook_name", "test-addrbook"),
					resource.TestCheckResourceAttr("ksyun_kfw_addrbook.foo", "ip_version", "IPv4"),
					resource.TestCheckResourceAttr("ksyun_kfw_addrbook.foo", "ip_address.#", "2"),
				),
			},
		},
	})
}

func TestAccKsyunCfwAddrbook_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_kfw_addrbook.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKfwAddrbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKfwAddrbookConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKfwAddrbookExists("ksyun_kfw_addrbook.foo", &val),
					testAccCheckKfwAddrbookAttributes(&val),
				),
			},
			{
				Config: testAccKfwAddrbookUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKfwAddrbookExists("ksyun_kfw_addrbook.foo", &val),
					testAccCheckKfwAddrbookAttributes(&val),
					resource.TestCheckResourceAttr("ksyun_kfw_addrbook.foo", "addrbook_name", "test-addrbook-updated"),
					resource.TestCheckResourceAttr("ksyun_kfw_addrbook.foo", "ip_address.#", "3"),
				),
			},
		},
	})
}

func testAccCheckKfwAddrbookExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return nil
}

func testAccCheckKfwAddrbookAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val == nil || len(*val) == 0 {
			return fmt.Errorf("CfwAddrbook attributes is empty")
		}
		return nil
	}
}

func testAccCheckKfwAddrbookDestroy(s *terraform.State) error {
	return nil
}

const testAccKfwAddrbookConfig = `
resource "ksyun_kfw_instance" "foo" {
  instance_name  = "test-Kfw-instance"
  instance_type  = "Advanced"
  bandwidth      = 50
  total_eip_num  = 50
  charge_type    = "Monthly"
  project_id     = "0"
  purchase_time  = 1
}

resource "ksyun_kfw_addrbook" "foo" {
  cfw_instance_id = ksyun_kfw_instance.foo._id
  addrbook_name   = "test-addrbook"
  ip_version      = "IPv4"
  ip_address      = ["10.1.1.11", "10.2.2.21"]
  description     = "test address book"
}
`

const testAccKfwAddrbookUpdateConfig = `
resource "ksyun_kfw_instance" "foo" {
  instance_name  = "test-kfw-instance"
  instance_type  = "Advanced"
  bandwidth      = 50
  total_eip_num  = 50
  charge_type    = "Monthly"
  project_id     = "0"
  purchase_time  = 1
}

resource "ksyun_kfw_addrbook" "foo" {
  kfw_instance_id = ksyun_kfw_instance.foo.id
  addrbook_name   = "test-addrbook-updated"
  ip_version      = "IPv4"
  ip_address      = ["10.1.1.11", "10.2.2.21", "10.3.3.31"]
  description     = "test address book updated"
}
`
