package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunSecurityGroupEntryLite_basic(t *testing.T) {
	// var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_security_group_entry_lite.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupEntryLiteConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccKsyunSecurityGroupEntryLite_update(t *testing.T) {
	// var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_security_group_entry_lite.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupEntryLiteConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccSecurityGroupEntryLiteUpdateConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccSecurityGroupEntryLiteConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="ksyun-security-group"
}

resource "ksyun_security_group_entry_lite" "foo" {
  description = "test1"
  security_group_id="${ksyun_security_group.default.id}"
  direction="in"
  protocol="ip"
  icmp_type=0
  icmp_code=0
  port_range_from=22
  port_range_to=22

  cidr_block=["192.177.2.1/32", "192.177.2.2/32", "10.0.1.3/32"]
}


resource "ksyun_security_group_entry_lite" "icmp" {
  description       = "test3"
  security_group_id = ksyun_security_group.default.id
  direction         = "in"
  protocol          = "tcp"
    # icmp_type         = 0
    # icmp_code         = 0
  port_range_from = 22
  port_range_to   = 22

  cidr_block = ["192.177.2.1/32", "192.177.2.2/32", "10.0.1.3/32"]
}

`

const testAccSecurityGroupEntryLiteUpdateConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="ksyun-security-group"
}
resource "ksyun_security_group_entry_lite" "foo" {
  description = "test2"
  security_group_id="${ksyun_security_group.default.id}"
  direction="in"
  protocol="ip"
  icmp_type=0
  icmp_code=0
  port_range_from=22
  port_range_to=22

  cidr_block=[ "192.177.2.1/32", "10.0.1.1/32"]
}
`
