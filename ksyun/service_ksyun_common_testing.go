package ksyun

import (
	"fmt"
	"strings"
)

func testBasicNetworkConfig(region, suffix string) string {
	s := fmt.Sprintf(`
provider "ksyun" {
  region = "%s"
}

resource "ksyun_vpc" "foo" {
  vpc_name   = "tf-${var.suffix}-vpc"
  cidr_block = "10.7.0.0/21"
}

data "ksyun_availability_zones" "foo" {
}

resource "ksyun_subnet" "foo" {
  subnet_name = "tf-${var.suffix}-subnet"
  cidr_block  = "10.7.0.0/21"
  subnet_type = "Normal"
  vpc_id      = ksyun_vpc.foo.id
  gateway_ip  = "10.7.0.1"
  dns1        = "198.18.254.41"
  dns2        = "198.18.254.40"
  #   provided_ipv6_cidr_block = true
  availability_zone = data.ksyun_availability_zones.foo.availability_zones.0.availability_zone_name
}

resource "ksyun_security_group" "foo" {
  vpc_id              = ksyun_vpc.foo.id
  security_group_name = "tf-${var.suffix}-security-group"
}

`, region)
	return strings.ReplaceAll(s, "${var.suffix}", suffix)

}
