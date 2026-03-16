package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunBareMetal_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_bare_metal.default",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckBareMetalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBareMetalConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBareMetalExists("ksyun_bare_metal.default", &val),
					testAccCheckBareMetalAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunBareMetal_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_bare_metal.default",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckBareMetalDestroy,
		Steps: []resource.TestStep{
			//{
			//	Config: testAccBareMetalConfig,
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheckBareMetalExists("ksyun_bare_metal.default", &val),
			//		testAccCheckBareMetalAttributes(&val),
			//	),
			//},
			{
				Config: testAccBareMetalUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBareMetalExists("ksyun_bare_metal.default", &val),
					testAccCheckBareMetalAttributes(&val),
					resource.TestCheckResourceAttr("ksyun_bare_metal.default", "host_name", "tf-test-0311-updated"),
				),
			},
		},
	})
}

func testAccCheckBareMetalExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Bare Metal id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		bareMetalService := BareMetalService{client}

		data, err := bareMetalService.ReadBareMetal(nil, rs.Primary.ID, true)
		if err != nil {
			return fmt.Errorf("error reading Bare Metal: %s", err)
		}

		if len(data) == 0 {
			return fmt.Errorf("Bare Metal not found")
		}

		*val = data
		return nil
	}
}

func testAccCheckBareMetalAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			if _, ok := (*val)["HostId"]; !ok {
				return fmt.Errorf("Bare Metal HostId not found")
			}
		}
		return nil
	}
}

func testAccCheckBareMetalDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_bare_metal" {
			continue
		}
		client := testAccProvider.Meta().(*KsyunClient)
		bareMetalService := BareMetalService{client}

		data, err := bareMetalService.ReadBareMetal(nil, rs.Primary.ID, true)

		if err == nil && data != nil && len(data) > 0 {
			return fmt.Errorf("Bare Metal still exist: %s", rs.Primary.ID)
		}
	}

	return nil
}

const testAccBareMetalConfig = `
resource "ksyun_bare_metal" "default" {
  availability_zone        = "cn-qingyangtest-1a"
  host_type                = "CAL"
  image_id                 = "460563e9-24f9-44e2-81c1-8143e8cb93c9"
  raid                     = "Raid1"
  network_interface_mode   = "single"
  subnet_id                = "835f125d-32c5-4e8f-a0ca-7474825d2bb0"
  security_group_ids       = ["7e94fede-52cb-40fd-82ac-99c9b21b5af4"]
  key_id                   = "6a1fc1a1-34d1-4c4a-bbce-3ecb8df855ef"
  host_name                = "tf-test-0311"
  charge_type              = "Daily"
  system_file_type         = "EXT4"
  custom_install_config {
      key = "grub_cmd_line"
      value = ["keepImagecredential"]
  }
  custom_install_config {
      key = "dracut_files"
      value = ["/etc/modprobe.d/test.conf","/etc/modprobe.d/nvidia.conf"]
  }
  force_re_install         = true
}
`

const testAccBareMetalUpdateConfig = `

resource "ksyun_bare_metal" "default" {
  availability_zone        = "cn-qingyangtest-1a"
  host_type                = "CAL"
  image_id                 = "460563e9-24f9-44e2-81c1-8143e8cb93c9"
  raid                     = "Raid1"
  network_interface_mode   = "single"
  subnet_id                = "835f125d-32c5-4e8f-a0ca-7474825d2bb0"
  security_group_ids       = ["7e94fede-52cb-40fd-82ac-99c9b21b5af4"]
  key_id                   = "6a1fc1a1-34d1-4c4a-bbce-3ecb8df855ef"
  host_name                = "tf-test-0311-updated"
  charge_type              = "Daily"
  system_file_type         = "EXT4"
  force_re_install         = true
  custom_install_config {
      key = "grub_cmd_line"
      value = ["keepImagecredential"]
  }
  custom_install_config {
      key = "dracut_files"
      value = ["/etc/modprobe.d/test.conf","/etc/modprobe.d/nvidia.conf"]
  }
}
`
