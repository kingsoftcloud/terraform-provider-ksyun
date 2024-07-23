package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunBackendServerGroupsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataBackendServerGroupsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_lb_backend_server_groups.foo"),
				),
			},
		},
	})
}

const testAccDataBackendServerGroupsConfig = `
data "ksyun_lb_backend_server_groups" "foo" {
  output_file="output_result"
  ids=[]
}
`
