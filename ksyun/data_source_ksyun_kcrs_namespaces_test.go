package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunKcrsNamespacesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKcrsNamespacesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kcrs_namespaces.foo"),
				),
			},
		},
	})
}

const testAccDataKcrsNamespacesConfig = `
data "ksyun_kcrs_namespaces" "foo" {
  output_file="kcrs_namespaces_output_result"
  instance_id = "86f14f8c-bf24-42c8-91bd-1a2c0dfd2224"
}
`
