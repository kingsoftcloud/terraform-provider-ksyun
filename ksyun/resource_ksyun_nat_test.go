package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunNat_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_nat.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNatDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccNatConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatExists("ksyun_nat.foo", &val),
					testAccCheckNatAttributes(&val),
					resource.TestCheckResourceAttr("ksyun_nat.foo", "nat_name", "ksyun-nat-tf"),
					resource.TestCheckResourceAttr("ksyun_nat.foo", "nat_mode", "Vpc"),
				),
			},
		},
	})
}

func TestAccKsyunNat_update(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_nat.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNatDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccNatConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatExists("ksyun_nat.foo", &val),
					testAccCheckNatAttributes(&val),
				),
			},
			{
				Config: testAccNatConfigUpdate,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatExists("ksyun_nat.foo", &val),
					testAccCheckNatAttributes(&val),
					resource.TestCheckResourceAttr("ksyun_nat.foo", "nat_name", "ksyun-nat-tf-update"),
				),
			},
		},
	})
}

func testAccCheckNatExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf(" Nat id is empty ")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		Nat := make(map[string]interface{})
		Nat["NatId.1"] = rs.Primary.ID
		ptr, err := client.vpcconn.DescribeNats(&Nat)

		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["NatSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}

func testAccCheckNatAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["NatSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf(" Nat id is empty ")
			}
		}
		return nil
	}
}

func testAccCheckNatDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_nat" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		Nat := make(map[string]interface{})
		Nat["NatId.1"] = rs.Primary.ID
		ptr, err := client.vpcconn.DescribeNats(&Nat)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["NatSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf(" Nat still exist ")
			}
		}
	}

	return nil
}

const testAccNatConfig = `
resource "ksyun_vpc" "test" {
  vpc_name = "ksyun-vpc-tf-test"
  cidr_block = "10.0.0.0/16"
}
resource "ksyun_nat" "foo" {
  nat_name = "ksyun-nat-tf-test"
  nat_mode = "Vpc"
  nat_type = "public"
  band_width = 1
  charge_type = "Peak"
  vpc_id = "${ksyun_vpc.test.id}"
}
`

const testAccNatConfigUpdate = `
resource "ksyun_vpc" "test" {
 vpc_name = "ksyun-vpc-tf"
 cidr_block = "10.0.0.0/16"
}
resource "ksyun_nat" "foo" {
 nat_name = "ksyun-nat-tf-test-update"
 nat_mode = "Vpc"
 nat_type = "public"
 band_width = 1
 charge_type = "Daily"
 vpc_id = "${ksyun_vpc.test.id}"
}
`
