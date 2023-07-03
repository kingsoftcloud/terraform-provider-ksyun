package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunRegisterBackendServersDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataRegisterBackendServersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_lb_register_backend_servers.foo"),
				),
			},
		},
	})
}

const testAccDataRegisterBackendServersConfig = `
data "ksyun_lb_register_backend_servers" "foo" {
  output_file="output_result"
  ids=[]
  backend_server_group_id=[]
}
`
