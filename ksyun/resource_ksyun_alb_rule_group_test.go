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
		},
	})
}

const testAccAlbRuleGroupConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}

# 监听器的转发策略
resource "ksyun_alb_rule_group" "default" {
  alb_listener_id         = "78769781-b887-4733-a58c-e89314c5925b"
  alb_rule_group_name     = "tf_alb_rule_group_unit_test"
  backend_server_group_id = "31bcee03-0401-4b6e-8a16-3b6765466b3c"
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

  listener_sync = "on"

}
`
