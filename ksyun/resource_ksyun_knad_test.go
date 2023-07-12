package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunKnad_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_knad.foo1",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKnadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKnadConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKnadExists("ksyun_knad.foo1", &val),
					testAccCheckKnadAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunKnad_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_knad.foo1",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKnadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKnadConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKnadExists("ksyun_knad.foo1", &val),
					testAccCheckKnadAttributes(&val),
				),
			},
			{
				Config: testAccKnadUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKnadExists("ksyun_knad.foo1", &val),
					testAccCheckKnadAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckKnadExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Knad id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		knad := make(map[string]interface{})
		knad["KnadId.1"] = rs.Primary.ID
		//knad["ProjectId.1"] = rs.Primary.Attributes["project_id"]
		ptr, err := client.knadconn.DescribeKnad(&knad)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["KnadSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}
		*val = *ptr
		return nil
	}
}
func testAccCheckKnadAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["KnadSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("Knad id is empty")
			}
		}
		return nil
	}
}
func testAccCheckKnadDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_knad" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		knad := make(map[string]interface{})
		knad["KnadId.1"] = rs.Primary.ID
		ptr, err := client.knadconn.DescribeKnad(&knad)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["KnadSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("Knad still exist")
			}
		}
	}

	return nil
}

const testAccKnadConfig = `

# Create an knad
resource "ksyun_knad" "foo1" {
  link_type = "DDoS_BGP"
  ip_count = 10
  band = 30
  max_band = 30
  idc_band = 100
  duration = 1
  knad_name = "test3"
  bill_type = 1
  service_id = "KNAD_30G"
  project_id="0"
}
`

const testAccKnadUpdateConfig = `

# Create an knad
resource "ksyun_knad" "foo1" {
   link_type = "DDoS_BGP"
  ip_count = 10
  band = 30
  max_band = 50
  idc_band = 100
  duration = 1
  knad_name = "test3"
  bill_type = 1
  service_id = "KNAD_30G"
  project_id="0"
}
`
