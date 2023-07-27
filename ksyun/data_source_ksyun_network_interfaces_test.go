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
  vpc_id=["9c1238d7-bd64-4e66-a784-13e5bd6e9a4b"]
  subnet_id=["73a8c6d4-d94c-4c26-8e01-5b3f9adea51a"]
  securitygroup_id=[]
  instance_type=[]
  instance_id=["0b49bce7-5c70-4416-a72c-5b69a284a7ce"]
  private_ip_address=[]
}
`
