package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunDnatsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDnatsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_dnats.foo"),
				),
			},
		},
	})
}

const testAccDataDnatsConfig = `
provider "ksyun" {
	region = "cn-guangzhou-1"
}
data "ksyun_dnats" "foo" {
  private_ip_address = "10.7.5.213"
  nat_id = "5c7b7925-a7d7-4db7-ac7a-434fc8042329"
  output_file = "output_result_dnats"
}
`
