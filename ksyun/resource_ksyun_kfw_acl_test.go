package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunCfwAcl_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_cfw_acl.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCfwAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwAclConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCfwAclExists("ksyun_cfw_acl.foo", &val),
					testAccCheckCfwAclAttributes(&val),
					resource.TestCheckResourceAttr("ksyun_cfw_acl.foo", "acl_name", "test-acl-rule"),
					resource.TestCheckResourceAttr("ksyun_cfw_acl.foo", "direction", "in"),
					resource.TestCheckResourceAttr("ksyun_cfw_acl.foo", "src_type", "ip"),
					resource.TestCheckResourceAttr("ksyun_cfw_acl.foo", "dest_type", "ip"),
					resource.TestCheckResourceAttr("ksyun_cfw_acl.foo", "service_type", "service"),
					resource.TestCheckResourceAttr("ksyun_cfw_acl.foo", "app_type", "any"),
					resource.TestCheckResourceAttr("ksyun_cfw_acl.foo", "policy", "accept"),
					resource.TestCheckResourceAttr("ksyun_cfw_acl.foo", "status", "start"),
				),
			},
		},
	})
}

func TestAccKsyunCfwAcl_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_cfw_acl.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCfwAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwAclConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCfwAclExists("ksyun_cfw_acl.foo", &val),
					testAccCheckCfwAclAttributes(&val),
				),
			},
			{
				Config: testAccCfwAclUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCfwAclExists("ksyun_cfw_acl.foo", &val),
					testAccCheckCfwAclAttributes(&val),
					resource.TestCheckResourceAttr("ksyun_cfw_acl.foo", "acl_name", "test-acl-rule-updated"),
					resource.TestCheckResourceAttr("ksyun_cfw_acl.foo", "policy", "deny"),
				),
			},
		},
	})
}

func testAccCheckCfwAclExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return nil
}

func testAccCheckCfwAclAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val == nil || len(*val) == 0 {
			return fmt.Errorf("CfwAcl attributes is empty")
		}
		return nil
	}
}

func testAccCheckCfwAclDestroy(s *terraform.State) error {
	return nil
}

const testAccCfwAclConfig = `
resource "ksyun_cfw_instance" "foo" {
  instance_name  = "test-cfw-instance"
  instance_type  = "Advanced"
  bandwidth      = 50
  total_eip_num  = 50
  charge_type    = "Monthly"
  project_id     = "0"
  purchase_time  = 1
}

resource "ksyun_cfw_acl" "foo" {
  cfw_instance_id   = ksyun_cfw_instance.foo.cfw_instance_id
  acl_name         = "test-acl-rule"
  direction        = "in"
  src_type         = "ip"
  src_ips          = ["1.1.1.1", "2.2.2.2"]
  dest_type        = "ip"
  dest_ips         = ["3.3.3.3"]
  service_type     = "service"
  service_infos     = ["TCP:1-100/80-80"]
  app_type         = "any"
  policy           = "accept"
  status           = "start"
  priority_position = "after+1"
  description      = "test acl rule"
}
`

const testAccCfwAclUpdateConfig = `
resource "ksyun_cfw_instance" "foo" {
  instance_name  = "test-cfw-instance"
  instance_type  = "Advanced"
  bandwidth      = 50
  total_eip_num  = 50
  charge_type    = "Monthly"
  project_id     = "0"
  purchase_time  = 1
}

resource "ksyun_cfw_acl" "foo" {
  cfw_instance_id   = ksyun_cfw_instance.foo.cfw_instance_id
  acl_name         = "test-acl-rule-updated"
  direction        = "in"
  src_type         = "ip"
  src_ips          = ["1.1.1.1", "2.2.2.2"]
  dest_type        = "ip"
  dest_ips         = ["3.3.3.3"]
  service_type     = "service"
  service_infos     = ["TCP:1-100/80-80"]
  app_type         = "any"
  policy           = "deny"
  status           = "start"
  priority_position = "after+1"
  description      = "test acl rule updated"
}
`
