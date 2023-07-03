package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunAvailabilityZonesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAvailabilityZonesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_availability_zones.foo"),
				),
			},
		},
	})
}

const testAccDataAvailabilityZonesConfig = `

data "ksyun_availability_zones" "foo" {
	output_file = "output_result"
}
`
