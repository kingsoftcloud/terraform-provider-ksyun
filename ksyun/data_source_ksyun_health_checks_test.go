package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunHealthChecksDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataHealthChecksConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_health_checks.foo"),
				),
			},
		},
	})
}

const testAccDataHealthChecksConfig = `
data "ksyun_health_checks" "foo" {
  output_file="output_result"
  #listener_id=["8d1dac22-6c6c-42ea-93e2-2702d44ddb93","70467f7e-23dc-465a-a609-fb1525fc6b16"]
}
`
