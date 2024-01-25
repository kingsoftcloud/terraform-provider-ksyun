package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunSecurityGroupEntrySet_basic(t *testing.T) {
	// var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_security_group_entry_Set.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupEntrySetConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccKsyunSecurityGroupEntrySet_update(t *testing.T) {
	// var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_security_group_entry_Set.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupEntrySetConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccSecurityGroupEntrySetUpdateConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccSecurityGroupEntrySetConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="ksyun-security-group"
}

resource "ksyun_security_group_entry_set" "foo" {
  security_group_id="${ksyun_security_group.default.id}"

  security_group_entries {
      description = "test1"
	  direction="in"
	  protocol="ip"
	  icmp_type=0
	  icmp_code=0
	  port_range_from=22
	  port_range_to=22
      cidr_block = "10.7.110.220/32"
  }
  security_group_entries {
      description = "test2"
	  direction="in"
	  protocol="ip"
	  icmp_type=0
	  icmp_code=0
	  port_range_from=22
	  port_range_to=22
      cidr_block = "10.7.111.220/32"
  }
}`

const testAccSecurityGroupEntrySetUpdateConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="ksyun-security-group"
}
resource "ksyun_security_group_entry_set" "foo" {
  security_group_id="${ksyun_security_group.default.id}"

  security_group_entries {
      description = "test1"
	  direction="in"
	  protocol="ip"
	  icmp_type=0
	  icmp_code=0
	  port_range_from=22
	  port_range_to=22
      cidr_block = "10.7.110.220/32"
  }
  security_group_entries {
      description = "test2"
	  direction="in"
	  protocol="ip"
	  icmp_type=0
	  icmp_code=0
	  port_range_from=22
	  port_range_to=22
      cidr_block = "10.7.111.220/32"
  }
}
`
