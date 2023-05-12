package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunMongodbsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataMongodbsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbExists("data.ksyun_mongodbs.default"),
				),
			},
		},
	})
}

func testAccCheckMongodbExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}
		if rs.Primary.Attributes["instances.#"] == "" {
			return fmt.Errorf("mongodb instance is not be set")
		}
		return nil
	}
}

const testAccDataMongodbsConfig = `
data "ksyun_mongodbs" "default" {
  output_file = "output_result"
  iam_project_id = ""
  instance_id = ""
  vnet_id = ""
  vpc_id = ""
  name = ""
  vip = ""
}
`
