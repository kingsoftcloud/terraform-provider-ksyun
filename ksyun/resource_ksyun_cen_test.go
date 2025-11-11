package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunCen_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_cen.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCenConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenExists("ksyun_cen.foo", &val),
					testAccCheckCenAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunCen_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_cen.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCenConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenExists("ksyun_cen.foo", &val),
					testAccCheckCenAttributes(&val),
				),
			},
			{
				Config: testAccCenUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenExists("ksyun_cen.foo", &val),
					testAccCheckCenAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckCenExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Cen id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		cen := make(map[string]interface{})
		cen["CenId.1"] = rs.Primary.ID
		ptr, err := client.cenconn.DescribeCens(&cen)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["CenSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}
		*val = *ptr
		return nil
	}
}
func testAccCheckCenAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["CenSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("Cen id is empty")
			}
		}
		return nil
	}
}
func testAccCheckCenDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_cen" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		cen := make(map[string]interface{})
		cen["CenId.1"] = rs.Primary.ID
		ptr, err := client.cenconn.DescribeCens(&cen)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["CenSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("Cen still exist")
			}
		}
	}

	return nil
}

const testAccCenConfig = `
resource "ksyun_cen" "foo" {
	cen_name="cen_create"
	description="zice_create"
}
`

const testAccCenUpdateConfig = `
resource "ksyun_cen" "foo" {
	cen_name="cen_update"
	description="zice_update"
}
`
