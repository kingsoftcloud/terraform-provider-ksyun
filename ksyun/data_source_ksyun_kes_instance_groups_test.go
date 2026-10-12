package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunKesInstanceGroups_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKesInstanceGroupsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kes_instance_groups.default"),
					resource.TestCheckResourceAttrSet("data.ksyun_kes_instance_groups.default", "instance_groups.0.id"),
				),
			},
		},
	})
}

const testAccDataKesInstanceGroupsConfig = `
provider "ksyun" {
  access_key = "your ak"
  secret_key = "your sk"
  region     = "cn-beijing-6"
  endpoint   = "kes.api.ksyun.com"
}

data "ksyun_kes_instance_groups" "default" {
  cluster_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  output_file = "output_result"
}

output "first_instance_group_id" {
  value = data.ksyun_kes_instance_groups.default.instance_groups.0.id
}
`
