package ksyun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKsyunAlbRegisterBackendServer_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_alb_register_backend_server.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlbRegisterBackendServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlbRegisterBackendServerConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccKsyunAlbRegisterBackendServer_update(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_alb_register_backend_server.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlbRegisterBackendServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlbRegisterBackendServerConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccAlbRegisterBackendServerUpdateConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccCheckAlbRegisterBackendServerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_alb_register_backend_server" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		backendServerGroup := make(map[string]interface{})
		backendServerGroup["BackendServerGroupId.1"] = rs.Primary.ID

		albSrv := AlbService{client: client}
		data, err := albSrv.ReadAlbBackendServer(nil, rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if notFoundError(err) {
				return nil
			} else {
				return err
			}
		}
		if data != nil {
			l := data
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("AlbBackendServer still exist")
			}
		}
	}

	return nil
}

const testAccAlbRegisterBackendServerConfig = `

resource "ksyun_alb_register_backend_server" "foo" {
  backend_server_group_id="0cc6ec96-0bf6-41b7-81e8-eba27af10f13"
  network_interface_id="fe84c574-be5a-43f9-803d-54bceee14411"
  backend_server_ip="10.5.0.171"
  port = 80
  weight=30
}
`

const testAccAlbRegisterBackendServerUpdateConfig = `
resource "ksyun_alb_register_backend_server" "foo" {
  backend_server_group_id="0cc6ec96-0bf6-41b7-81e8-eba27af10f13"
  network_interface_id="fe84c574-be5a-43f9-803d-54bceee14411"
  backend_server_ip="10.5.0.171"
  port = 8080
  weight=40
}
`
