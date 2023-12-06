package ksyun

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func Test_resourceKsyunVpnGatewayRoute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_vpn_gateway_route.default",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVPCDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccVpnGatewayRouteConfig("vpn-route-unit-test"),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_vpn_gateway_route.default"),
				),
			},
		},
	})
}

func testAccVpnGatewayRouteConfig(suffix string) (s string) {
	defer func() {
		s = strings.ReplaceAll(s, "${var.suffix}", suffix)
	}()
	basicConfig := testBasicNetworkConfig("cn-guangzhou-1", suffix)
	return fmt.Sprintf(`
	%s
resource "ksyun_vpn_gateway" "default" {
  vpn_gateway_name   = "tf-${var.suffix}-vpn-gw"
  band_width = 10
  vpc_id = ksyun_vpc.foo.id
  charge_type = "Daily"
  vpn_gateway_version = "2.0"
}

resource "ksyun_vpc" "foo1" {
	vpc_name        = "tf-${var.suffix}-vpc"
	cidr_block = "10.7.248.0/21"
}

resource "ksyun_vpn_gateway_route" "default" {
  vpn_gateway_id = ksyun_vpn_gateway.default.id
  next_hop_type = "vpc"
  // next_hop_instance_id = ksyun_vpc.foo1.id
  destination_cidr_block = "10.7.248.0/21"
}

`, basicConfig)
}
