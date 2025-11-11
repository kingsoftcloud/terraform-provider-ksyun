package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunKfwAddrbooksDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKfwAddrbooksConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kfw_addrbooks.default"),
				),
			},
		},
	})
}

func TestAccKsyunKfwAddrbooksDataSource_withFilters(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKfwAddrbooksConfigWithFilters,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kfw_addrbooks.filtered"),
				),
			},
		},
	})
}

const testAccDataKfwAddrbooksConfig = `
data "ksyun_kfw_addrbooks" "default" {
  output_file       = "output_result"
  cfw_instance_id   = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  ids               = []
}
`

const testAccDataKfwAddrbooksConfigWithFilters = `
data "ksyun_kfw_addrbooks" "filtered" {
  output_file       = "output_result_filtered"
  cfw_instance_id   = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  ids               = ["6e218619-c681-4b9b-bef4-0cbb955b103e", "4ba903ab-e71c-4f88-9433-6db28dee9161"]
}
`
