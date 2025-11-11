package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunKfwInstancesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKfwInstancesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kfw_instances.default"),
				),
			},
		},
	})
}

func TestAccKsyunKfwInstancesDataSource_withFilters(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKfwInstancesConfigWithFilters,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kfw_instances.filtered"),
				),
			},
		},
	})
}

const testAccDataKfwInstancesConfig = `
data "ksyun_kfw_instances" "default" {
  output_file = "output_result"
  ids         = []
}
`

const testAccDataKfwInstancesConfigWithFilters = `
data "ksyun_kfw_instances" "filtered" {
  output_file = "output_result_filtered"
  ids         = ["fe56b633-cde1-4f96-8573-30378d1406b1", "dc028867-4a96-4982-95ca-99dbbe96d48a"]
}
`
