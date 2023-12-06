package ksyun

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunVpnTunnel_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_vpn_tunnel.default",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVPCDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccVpnTunnelConfig("vpn-tunnel-unit-test"),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_vpn_tunnel.default"),
				),
			},
		},
	})
}

func testAccVpnTunnelConfig(suffix string) (s string) {
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

resource "ksyun_vpn_customer_gateway" "default" {
  customer_gateway_address   = "100.0.0.65"
  ha_customer_gateway_address = "100.0.2.65"
  customer_gateway_name = "tf-${var.suffix}-vpn-cgw"
}

resource "ksyun_vpn_tunnel" "default" {
  vpn_tunnel_name = "tf-${var.suffix}-vpn-tunnel"
  type = "RouteIpsec"
  vpn_gateway_version = "2.0"
  vpn_gateway_id = ksyun_vpn_gateway.default.id
  customer_gateway_id = ksyun_vpn_customer_gateway.default.id
  ike_dh_group = 2
  ike_version = "v2"
  pre_shared_key = "123456789abcd"
  customer_peer_ip = "2.2.2.1/30"
  local_peer_ip = "2.2.2.2/30"
}

`, basicConfig)
}
