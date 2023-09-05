package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestResourceKsyunAlbRuleGroup_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// IDRefreshName: "ksyun_auto_snapshot_volume_association.foo",
		Providers: testAccProviders,
		// CheckDestroy:  testAccCheckSnapshotDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccAlbRuleGroupConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_alb_rule_group.default"),
				),
			},
			{
				Config: testAccAlbRuleGroupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_alb_rule_group.default"),
				),
			},
		},
	})
}

const testAccAlbRuleGroupConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}

# 监听器的转发策略
resource "ksyun_alb_rule_group" "default" {
  alb_listener_id         = "abf046dd-ce3f-4847-a6f4-dab413d222fe"
  alb_rule_group_name     = "tf_alb_rule_group_unit_test"
  backend_server_group_id = "7493cfee-afd4-4db4-af56-96b17491cfb6"
  alb_rule_set {
    # domain = "www.ksyun.com"
    # url = "/test/path"
    alb_rule_type = "url"
    alb_rule_value = "/test/path"
  }
alb_rule_set {
    # domain = "www.ksyun.com"
    alb_rule_type = "domain"
    alb_rule_value = "www.ksyun.com"
  }

  listener_sync = "off"
  session_state              = "start"
  session_persistence_period = 333
  cookie_type                = "ImplantCookie"
  # cookie_name                = "dasdad"
  health_check_state         = "start"
  interval = 3
  timeout=4
  healthy_threshold=3
  unhealthy_threshold=5
  host_name = ""
}
`
const testAccAlbRuleGroupUpdateConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}

# 监听器的转发策略
resource "ksyun_alb_rule_group" "default" {
  alb_listener_id         = "abf046dd-ce3f-4847-a6f4-dab413d222fe"
  alb_rule_group_name     = "tf_alb_rule_group_unit_test"
  backend_server_group_id = "7493cfee-afd4-4db4-af56-96b17491cfb6"
  alb_rule_set {
    # domain = "www.ksyun.com"
    # url = "/test/path"
    alb_rule_type = "url"
    alb_rule_value = "/test/path"
  }
alb_rule_set {
    # domain = "www.ksyun.com"
    alb_rule_type = "domain"
    alb_rule_value = "www.ksyun.com"
  }

  listener_sync = "off"
  session_state              = "start"
  session_persistence_period = 333
  cookie_type                = "RewriteCookie"
  # cookie_name                = "dasdad"
  health_check_state         = "start"
  interval = 3
  timeout=4
  healthy_threshold=3
  unhealthy_threshold=5
  host_name = ""
}
`
