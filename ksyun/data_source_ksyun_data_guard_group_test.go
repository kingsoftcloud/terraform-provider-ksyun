package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunDataGuardGroupDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDataGuardGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_data_guard_group.foo"),
				),
			},
		},
	})
}

const testAccDataDataGuardGroupConfig = `
provider "ksyun" {
	region = "cn-qingyangtest-1"
}


data "ksyun_data_guard_group" "foo" {
  output_file = "output_result"
}
`
