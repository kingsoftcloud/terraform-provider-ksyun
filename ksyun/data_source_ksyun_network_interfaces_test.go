package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunNetworkInterfacesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataNetworkInterfacesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_network_interfaces.foo"),
				),
			},
		},
	})
}

const testAccDataNetworkInterfacesConfig = `

provider "ksyun" {
	region = "cn-guangzhou-1"
}

data "ksyun_network_interfaces" "foo" {
  output_file="output_result_kni"
  ids=["30c5e86f-d938-490c-b5b4-fbf54cc1f5d4"]
  vpc_id=[]
  subnet_id=[]
  securitygroup_id=[]
  instance_type=[]
  instance_id=[]
  private_ip_address=[]
}
`
