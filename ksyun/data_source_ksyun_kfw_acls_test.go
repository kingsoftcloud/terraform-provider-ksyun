package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunKfwAclsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKfwAclsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kfw_acls.default"),
				),
			},
		},
	})
}

func TestAccKsyunKfwAclsDataSource_withFilters(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKfwAclsConfigWithFilters,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kfw_acls.filtered"),
				),
			},
		},
	})
}

const testAccDataKfwAclsConfig = `
data "ksyun_kfw_acls" "default" {
  output_file       = "output_result"
  cfw_instance_id   = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  ids               = []
}
`

const testAccDataKfwAclsConfigWithFilters = `
data "ksyun_kfw_acls" "filtered" {
  output_file       = "output_result_filtered"
  cfw_instance_id   = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  ids               = ["07d0fc33-00be-45b6-a167-e20c4e954b7d", "20a94c10-d6d0-417c-838e-d0dde0f860ca"]
}
`
