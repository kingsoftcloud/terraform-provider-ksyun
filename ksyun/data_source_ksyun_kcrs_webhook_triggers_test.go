package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunKcrsWebhookTriggersDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKcrsWebhookTriggersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kcrs_webhook_triggers.foo"),
				),
			},
		},
	})
}

const testAccDataKcrsWebhookTriggersConfig = `
data "ksyun_kcrs_webhook_triggers" "foo" {
  output_file="kcrs_webhook_triggers_output_result"
  instance_id = "86f14f8c-bf24-42c8-91bd-1a2c0dfd2224"
  namespace = "tftest"
}
`
