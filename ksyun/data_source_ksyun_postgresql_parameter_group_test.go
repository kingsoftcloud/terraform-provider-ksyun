package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunPostgresqlParameterGroupDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// {
			// 	Config: testAccDataPostgresqlParameterGroupConfig,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckIDExists("data.ksyun_postgresql_parameter_group.foo"),
			// 	),
			// },
			{
				Config: testAccDataPostgresqlParameterGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_postgresql_parameter_group.foo"),
				),
			},
		},
	})
}

const testAccDataPostgresqlParameterGroupConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}


data "ksyun_postgresql_parameter_group" "foo" {
	output_file = "output_result"
	db_parameter_group_id = "e217a352-6262-47ac-8bc6-dfd9551e7643"
}
`
const testAccDataPostgresqlParameterGroupKeywordConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}


data "ksyun_postgresql_parameter_group" "bra" {
	output_file = "output_result_keyword"
	// db_parameter_group_id = "b233609c-42e1-4aad-aa68-9a2ebdf68a82"
	keyword = "tf"
}
`
