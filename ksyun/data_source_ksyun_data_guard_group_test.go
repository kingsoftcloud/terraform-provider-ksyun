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
	region = "cn-beijing-6"
}

data "ksyun_data_guard_group" "foo" {
  	output_file = "output_result_dgp"
}
data "ksyun_data_guard_group" "foo1" {
	output_file = "output_result_dgp_name"
	data_guard_name = "dataGuardNameTEST"
}

data "ksyun_data_guard_group" "foo2" {
	output_file = "output_result_dgp_id"
	data_guard_id = "36323f6c-4c05-4739-937c-a181c92c29bb"
}

`
