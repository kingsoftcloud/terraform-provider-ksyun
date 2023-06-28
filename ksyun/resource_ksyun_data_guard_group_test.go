package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestResourceKsyunDataGuardGroup_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_data_guard_group.foo",
		Providers:     testAccProviders,
		// CheckDestroy:  testAccCheckSnapshotDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccDataGuardGroupConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_data_guard_group.foo"),
				),
			},
			// // to test terraform when its configuration changes
			{
				Config: testAccDataGuardGroupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_data_guard_group.foo"),
				),
			},
		},
	})
}

const testAccDataGuardGroupConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}


resource "ksyun_data_guard_group" "foo" {
  data_guard_name = "terraform_test_on_hcl"
  data_guard_type = "host"
}
`

const testAccDataGuardGroupUpdateConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}


resource "ksyun_data_guard_group" "foo" {
  data_guard_name = "terraform_test_on_hcl_rename"
  data_guard_type = "host"
}
`
