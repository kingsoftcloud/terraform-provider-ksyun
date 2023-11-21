package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunPdnsZoneDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPdnsZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_private_dns_zones.foo"),
				),
			},
		},
	})
}

const testAccDataPdnsZoneConfig = `

data "ksyun_private_dns_zones" "foo" {
  output_file = "pdns_output_result"
  zone_ids = []
}
`
