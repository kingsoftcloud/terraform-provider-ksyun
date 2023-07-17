package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunKnadsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKnadsConfig,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["data.ksyun_knads.foo"]

						if !ok {
							return fmt.Errorf(" Can't find resource or data source: %s ", "data.ksyun_knads.foo")
						}

						if rs.Primary.ID == "" {
							return fmt.Errorf("ID is not be set")
						}
						return nil
					}),
			},
		},
	})
}

const testAccDataKnadsConfig = `

data "ksyun_knads" "foo" {
 ids = ["xxxx-xxxxx"]
project_id = []

}
`
