package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunScalingConfiguration_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_scaling_configuration.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckScalingConfigurationDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccScalingConfigurationConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingConfigurationExists("ksyun_scaling_configuration.foo", &val),
					testAccCheckScalingConfigurationAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunScalingConfiguration_update(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_scaling_configuration.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckScalingConfigurationDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccScalingConfigurationConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingConfigurationExists("ksyun_scaling_configuration.foo", &val),
					testAccCheckScalingConfigurationAttributes(&val),
				),
			},
			{
				Config: testAccScalingConfigurationConfigUpdate,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingConfigurationExists("ksyun_scaling_configuration.foo", &val),
					testAccCheckScalingConfigurationAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunScalingConfiguration_essd(t *testing.T) {
	var val map[string]interface{}
	r := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_scaling_configuration.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckScalingConfigurationDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccScalingConfigurationConfigEssd(r),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingConfigurationExists("ksyun_scaling_configuration.foo", &val),
					resource.TestCheckResourceAttr(
						"ksyun_scaling_configuration.foo", "system_disk_type", "ESSD_SYSTEM_PL1"),
					resource.TestCheckResourceAttr(
						"ksyun_scaling_configuration.foo", "data_disks.0.disk_type", "ESSD_PL1"),
					resource.TestCheckResourceAttr(
						"ksyun_scaling_configuration.foo", "data_disks.0.disk_size", "50"),
				),
			},
		},
	})
}

func testAccCheckScalingConfigurationExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf(" ScalingConfiguration id is empty ")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		req := make(map[string]interface{})
		req["ScalingConfigurationId.1"] = rs.Primary.ID
		ptr, err := client.kecconn.DescribeScalingConfiguration(&req)

		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["ScalingConfigurationSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}

func testAccCheckScalingConfigurationAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["ScalingConfigurationSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf(" ScalingConfiguration id is empty ")
			}
		}
		return nil
	}
}

func testAccCheckScalingConfigurationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_scaling_configuration" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		req := make(map[string]interface{})
		req["ScalingConfigurationId.1"] = rs.Primary.ID
		ptr, err := client.kecconn.DescribeScalingConfiguration(&req)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil && (*ptr)["ScalingConfigurationSet"] != nil {
			l := (*ptr)["ScalingConfigurationSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf(" ScalingConfiguration still exist ")
			}
		}
	}

	return nil
}

const testAccScalingConfigurationConfig = `
resource "ksyun_scaling_configuration" "foo" {
  scaling_configuration_name = "tf-xym-test-1"
  image_id = "IMG-5465174a-6d71-4770-b8e1-917a0dd92466"
  instance_type = "N3.1B"
  password = "Aa123456"
}
`

const testAccScalingConfigurationConfigUpdate = `
resource "ksyun_scaling_configuration" "foo" {
  scaling_configuration_name = "tf-xym-test-2"
  image_id = "IMG-5465174a-6d71-4770-b8e1-917a0dd92466"
  instance_type = "N3.1B"
  password = "Aa123456"
}
`

func testAccScalingConfigurationConfigEssd(r int) string {
	return fmt.Sprintf(`
resource "ksyun_scaling_configuration" "foo" {
  scaling_configuration_name = "tf-essd-test-%d"
  image_id                   = "IMG-5465174a-6d71-4770-b8e1-917a0dd92466"
  instance_type              = "S6.2B"
  password                   = "Aa123456"
  system_disk_type           = "ESSD_SYSTEM_PL1"
  system_disk_size           = 50

  data_disks {
    disk_type = "ESSD_PL1"
    disk_size = 50
  }
}
`, r)
}
