package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunKcrsWebhookTrigger_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_kcrs_webhook_trigger.foo",

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKsyunKcrsWebhookTriggerConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_kcrs_webhook_trigger.foo"),
				),
			},
		},
	})
}

const testAccKsyunKcrsWebhookTriggerConfig = `
resource "ksyun_kcrs_webhook_trigger" "foo" {
	instance_id = "b061ebeb-106d-40b9-88ea-7cad7e0c08e5"
	namespace = "lusongke"
	trigger {
		trigger_url = "http://www.test111.com"
		trigger_name = "tfunittest"
		event_types = ["DeleteImage", "PushImage"]
		headers {
			key = "pp1"
			value = "22"
		}
		headers {
			key = "pp1"
			value = "333"
		}
	}
}
`
