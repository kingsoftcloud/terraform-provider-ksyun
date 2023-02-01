package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunRabbitmqInstancesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataRabbitmqInstancesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRabbitmqExists("data.ksyun_rabbitmqs.default"),
				),
			},
		},
	})
}

func testAccCheckRabbitmqExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}
		if rs.Primary.Attributes["instances.#"] == "" {
			return fmt.Errorf("rabbitmq instance is not be set")
		}
		return nil
	}
}

const testAccDataRabbitmqInstancesConfig = `
data "ksyun_rabbitmqs" "default" {
  output_file = ""
  project_id            = 103800
  instance_id = ""
  instance_name = ""
  subnet_id = ""
  vpc_id = ""
  vip = ""
}
`
