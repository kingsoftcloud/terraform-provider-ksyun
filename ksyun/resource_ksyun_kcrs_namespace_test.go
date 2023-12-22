package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunKcrsNamespace_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_kcrs_namespace.foo",

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKsyunKcrsNamespaceConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_kcrs_namespace.foo"),
				),
			}, {
				Config: testAccKsyunKcrsNamespaceUpdateConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_kcrs_namespace.foo"),
				),
			},
		},
	})
}

const testAccKsyunKcrsNamespaceConfig = `
resource "ksyun_kcrs_namespace" "foo" {
	instance_id = "b061ebeb-106d-40b9-88ea-7cad7e0c08e5"
	namespace = "tftest"
	public = false
}
`
const testAccKsyunKcrsNamespaceUpdateConfig = `
resource "ksyun_kcrs_namespace" "foo" {
	instance_id = "b061ebeb-106d-40b9-88ea-7cad7e0c08e5"
	namespace = "tftest"
	public = true 
}
`
