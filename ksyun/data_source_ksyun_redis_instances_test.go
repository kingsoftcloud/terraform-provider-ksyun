package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunKcsDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKcsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKcsExists("data.ksyun_redis_instances.default"),
				),
			},
		},
	})
}

func testAccCheckKcsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}
		if rs.Primary.Attributes["instances.#"] == "" {
			return fmt.Errorf("kcs instance is not be set")
		}
		return nil
	}
}

const testAccDataKcsConfig = `
data "ksyun_redis_instances" "default" {
  output_file       = "output_result1"
  fuzzy_search      = ""
  iam_project_id    = ""
  cache_id          = ""
  vnet_id           = ""
  vpc_id            = ""
  name              = ""
  vip               = ""
}`
