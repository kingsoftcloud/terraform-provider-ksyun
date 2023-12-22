package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunKcrsToken_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_kcrs_token.foo",

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKsyunKcrsTokenConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_kcrs_token.foo"),
				),
			},
			{
				Config: testAccKsyunKcrsTokenUpdateConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_kcrs_token.foo"),
				),
			},
		},
	})
}

const testAccKsyunKcrsTokenConfig = `
resource "ksyun_kcrs_token" "foo" {
	instance_id = "b061ebeb-106d-40b9-88ea-7cad7e0c08e5"
	token_type = "Day"
	token_time = 10
	desc = "test"
	enable = true
}
`
const testAccKsyunKcrsTokenUpdateConfig = `
resource "ksyun_kcrs_token" "foo" {
	instance_id = "b061ebeb-106d-40b9-88ea-7cad7e0c08e5"
	token_type = "Day"
	token_time = 100
	desc = "test-11"
	enable = false
}
`
