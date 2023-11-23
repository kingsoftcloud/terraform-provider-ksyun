package ksyun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKsyunPrivateDnsZone_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_private_dns_zone.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckPrivateDnsDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateDnsExists("ksyun_private_dns_zone.foo", &val),
					testAccCheckPrivateDnsAttributes(&val),
				),
			},
		},
	})
}
func TestAccKsyunPrivateDnsZone_update(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_private_dns_zone.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckPrivateDnsDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateDnsExists("ksyun_private_dns_zone.foo", &val),
					testAccCheckPrivateDnsAttributes(&val),
				),
			},
			{
				Config: testAccPrivateDnsConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateDnsExists("ksyun_private_dns_zone.foo", &val),
					testAccCheckPrivateDnsAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckPrivateDnsExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("PrivateDnsZone id is empty")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		PrivateDns := make(map[string]interface{})
		PrivateDns["Filter.1"] = rs.Primary.ID
		ptr, err := client.pdnsconn.DescribePdnsZones(&PrivateDns)

		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["ZoneSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}

func testAccCheckPrivateDnsAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["ZoneSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("PrivateDns id is empty")
			}
		}
		return nil
	}
}

func testAccCheckPrivateDnsDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_private_dns_zone" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		PrivateDns := make(map[string]interface{})
		PrivateDns["Filter.1"] = rs.Primary.ID
		ptr, err := client.pdnsconn.DescribePdnsZones(&PrivateDns)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["ZoneSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("PrivateDns still exist")
			}
		}
	}

	return nil
}

const testAccPrivateDnsConfig = `
resource "ksyun_private_dns_zone" "foo" {
	zone_name = "tf-pdns-zone-pdns.com"
	zone_ttl = 360
	charge_type = "TrafficMonthly"
}
`

const testAccPrivateDnsConfigUpdate = `
resource "ksyun_private_dns_zone" "foo" {
	zone_name = "tf-pdns-zone-pdns.com"
	zone_ttl = 3999
	charge_type = "TrafficMonthly"
}
`
