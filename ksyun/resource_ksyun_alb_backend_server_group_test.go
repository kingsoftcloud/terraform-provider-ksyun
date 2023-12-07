package ksyun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKsyunAlbBackendServerGroup_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_alb_backend_server_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlbBackendServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlbBackendServerGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbBackendServerGroupExists("ksyun_alb_backend_server_group.foo", &val),
					testAccCheckAlbBackendServerGroupAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunAlbBackendServerGroup_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_alb_backend_server_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlbBackendServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlbBackendServerGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbBackendServerGroupExists("ksyun_alb_backend_server_group.foo", &val),
					testAccCheckAlbBackendServerGroupAttributes(&val),
				),
			},
			{
				Config: testAccAlbBackendServerGroupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbBackendServerGroupExists("ksyun_alb_backend_server_group.foo", &val),
					testAccCheckAlbBackendServerGroupAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckAlbBackendServerGroupExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("LbBackendServerGroup id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		backendServerGroup := make(map[string]interface{})
		backendServerGroup["BackendServerGroupId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeAlbBackendServerGroups(&backendServerGroup)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["BackendServerGroupSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}
		*val = *ptr
		return nil
	}
}
func testAccCheckAlbBackendServerGroupAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["BackendServerGroupSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("BackendServerGroup id is empty")
			}
		}
		return nil
	}
}
func testAccCheckAlbBackendServerGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_lb_backend_server_group" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		backendServerGroup := make(map[string]interface{})
		backendServerGroup["BackendServerGroupId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeBackendServerGroups(&backendServerGroup)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["BackendServerGroupSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("BackendServerGroup still exist")
			}
		}
	}

	return nil
}

const testAccAlbBackendServerGroupConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun_vpc_tf"
  cidr_block = "10.5.0.0/21"
}
resource "ksyun_alb_backend_server_group" "foo" {
  name="tf-alb-bsg"
  vpc_id="${ksyun_vpc.default.id}"
  upstream_keepalive="adaptation"
  backend_server_type="Host"
}
`

const testAccAlbBackendServerGroupUpdateConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun_vpc_tf"
  cidr_block = "10.5.0.0/21"
}
resource "ksyun_alb_backend_server_group" "foo" {
  name="tf-alb-bsg-update"
  vpc_id="${ksyun_vpc.default.id}"
  upstream_keepalive="adaptation"
  backend_server_type="Host"
}
`
