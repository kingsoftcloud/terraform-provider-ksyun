package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunKfwServiceGroupsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKfwServiceGroupsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kfw_service_groups.default"),
				),
			},
		},
	})
}

func TestAccKsyunKfwServiceGroupsDataSource_withFilters(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKfwServiceGroupsConfigWithFilters,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kfw_service_groups.filtered"),
				),
			},
		},
	})
}

const testAccDataKfwServiceGroupsConfig = `
data "ksyun_kfw_service_groups" "default" {
  output_file     = "output_result"
  cfw_instance_id = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  ids             = []
}
`

const testAccDataKfwServiceGroupsConfigWithFilters = `
data "ksyun_kfw_service_groups" "filtered" {
  output_file     = "output_result_filtered"
  cfw_instance_id = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  ids             = ["service-group-12345", "service-group-67890"]
}
`
