package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunKcrsInstancesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKcrsInstancesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kcrs_instances.foo"),
				),
			},
		},
	})
}

const testAccDataKcrsInstancesConfig = `
data "ksyun_kcrs_instances" "foo" {
  output_file="kcrs_instance_output_result"
  ids = []
}
`
