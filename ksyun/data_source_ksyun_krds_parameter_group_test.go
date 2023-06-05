package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunKrdsParameterGroupDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKrdsParameterGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_krds_parameter_group.foo"),
				),
			},
		},
	})
}

const testAccDataKrdsParameterGroupConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}


data "ksyun_krds_parameter_group" "foo" {
	output_file = "output_result"
	db_parameter_group_id = "b233609c-42e1-4aad-aa68-9a2ebdf68a82"
}
`
