package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunKnadAssociationAssociation_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_knad_associate.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKnadAssociationAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKnadAssociationAssociationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKnadAssociationAssociationExists("ksyun_knad_associate.foo", &val),
					testAccCheckKnadAssociationAssociationAttributes(&val),
				),
			},
			{
				Config: testAccKnadAssociationAssociationUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKnadAssociationAssociationExists("ksyun_knad_associate.foo", &val),
					testAccCheckKnadAssociationAssociationAttributes(&val),
				),
			},
		},
	})
}
func testAccCheckKnadAssociationAssociationExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("KnadAssociationAssociation id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		knadAssociationAssociation := make(map[string]interface{})
		knadAssociationAssociation["KnadId"] = rs.Primary.ID
		ptr, err := client.knadconn.IpList(&knadAssociationAssociation)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["EipSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("knad has no associate ips")
			}
		}

		*val = *ptr
		return nil
	}
}
func testAccCheckKnadAssociationAssociationAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["EipSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("KnadAssociationAssociation id is empty")
			}
		}
		return nil
	}
}
func testAccCheckKnadAssociationAssociationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_knad_associate" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		knadAssociationAssociation := make(map[string]interface{})
		knadAssociationAssociation["KnadId"] = rs.Primary.ID
		ptr, err := client.knadconn.IpList(&knadAssociationAssociation)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["EipSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("KnadAssociationAssociation still exist")
			}
		}
	}
	return nil
}

const testAccKnadAssociationAssociationConfig = `

resource "ksyun_knad_associate" "foo" {
  knad_id="xxx-xxx-xxx"
  ip = ["1.1.1.2"]

}
`
const testAccKnadAssociationAssociationUpdateConfig = `

resource "ksyun_knad_associate" "foo" {
  knad_id="xxxx-xxxx-xxx"
  ip = ["1.1.1.1"]

}
`
