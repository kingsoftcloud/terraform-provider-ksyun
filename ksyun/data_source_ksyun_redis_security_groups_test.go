package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunKcsSecGroupDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKcsSecGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKcsSecGroupExists("data.ksyun_redis_security_groups.default"),
				),
			},
		},
	})
}

func testAccCheckKcsSecGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}
		if rs.Primary.Attributes["instances.#"] == "" {
			return fmt.Errorf("kcs sec group is not be set")
		}
		return nil
	}
}

const testAccDataKcsSecGroupConfig = `
data "ksyun_redis_security_groups" "default" {
  output_file       = "output_result1"
}`
