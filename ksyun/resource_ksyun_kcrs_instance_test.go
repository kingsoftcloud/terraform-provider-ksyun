package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunKcrsInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKsyunKcrsInstanceConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_kcrs_instance.foo"),
				),
			},
			{
				Config: testAccKsyunKcrsInstanceUpdateConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_kcrs_instance.foo"),
				),
			},
		},
	})
}

const testAccKsyunKcrsInstanceConfig = `
resource "ksyun_kcrs_instance" "foo" {
	instance_name = "tfunittest1"
	instance_type = "basic"
	open_public_operation = true

	external_policy {
		entry = "192.168.2.133"
		desc = "ddd"
	}
	external_policy {
		entry = "192.168.2.123/32"
		desc = "ddd"
	}
}
`

const testAccKsyunKcrsInstanceUpdateConfig = `
resource "ksyun_kcrs_instance" "foo" {
	instance_name = "tfunittest1"
	instance_type = "basic"
	open_public_operation = true

	external_policy {
		entry = "192.168.2.133"
		desc = "ddd"
	}
	external_policy {
		entry = "192.168.2.122/32"
		desc = "ddd"
	}
}
`
