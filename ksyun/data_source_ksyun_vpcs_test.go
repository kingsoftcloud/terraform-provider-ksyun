package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunVPCsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataVPCsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_vpcs.foo"),
				),
			},
		},
	})
}

const testAccDataVPCsConfig = `
provider "ksyun" {
	region = "cn-guangzhou-1"
}
// resource "ksyun_vpc" "default" {
// 	vpc_name        = "tf-acc-vpc-data"
//     cidr_block      = "192.168.0.0/16"
// }
data "ksyun_vpcs" "foo" {
    ids = ["8155753a-4f3a-44fc-bef4-23a6e5aa4ad4"]
	output_file = "output_result_vpc"
}
`
