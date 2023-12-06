package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func Test_dataSourceKsyunVpnGatewayRoutes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataVpnGatewayRoutesConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_vpn_gateway_routes.foo"),
				),
			},
		},
	})
}

const testAccDataVpnGatewayRoutesConfig = `
provider "ksyun" {
	region = "cn-guangzhou-1"
}

data "ksyun_vpn_gateway_routes" "foo" {
	output_file = "output_result_vpn_gateway_routes_cidr"
	cidr_blocks = ["10.7.248.0/21"]
}

data "ksyun_vpn_gateway_routes" "foo1" {
	output_file = "output_result_vpn_gateway_routes_type"
	next_hop_types = ["vpc"]
}
`
