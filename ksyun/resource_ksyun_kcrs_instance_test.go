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
		},
	})
}

const testAccKsyunKcrsInstanceConfig = `
resource "ksyun_kcrs_instance" "foo" {
	instance_name = "tfunittest1"
	instance_type = "basic"
}
`
