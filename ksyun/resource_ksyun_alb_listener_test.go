package ksyun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKsyunAlbListener_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_alb_listener.test",
		Providers:     testAccProviders,
		// CheckDestroy:  testAccCheckAlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlbListenerConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccKsyunAlbListener_update(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_alb_listener.test",
		Providers:     testAccProviders,
		// CheckDestroy:  testAccCheckAlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlbListenerConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccAlbListenerUpdateConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccCheckAlbListenerDestroy(s *terraform.State) error {
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

const testAccAlbListenerConfig = `
resource "ksyun_alb_listener" "test" {
  alb_id             = "7291ecea-dcc5-45c1-8331-e53e062d2f52"
  alb_listener_name  = "alb-unit-test-listener"
  protocol           = "HTTPS"
  port               = 8087
  alb_listener_state = "start"
  certificate_id     = "01e82ae3-945a-4aad-a2af-e3c38eeea835"
  session {
    cookie_type                = "ImplantCookie"
    cookie_name                = "KLBRSIDdad"
    session_state              = "start"
    session_persistence_period = 3100
  }
  default_forward_rule {
    backend_server_group_id = "45390ce2-6703-4a4a-a662-9ca338b77c09"
  }
}
`

const testAccAlbListenerUpdateConfig = `
resource "ksyun_alb_listener" "test" {
  alb_id             = "7291ecea-dcc5-45c1-8331-e53e062d2f52"
  alb_listener_name  = "alb-unit-test-listener"
  protocol           = "HTTPS"
  port               = 8087
  alb_listener_state = "start"
  certificate_id     = "01e82ae3-945a-4aad-a2af-e3c38eeea835"
  session {
    cookie_type                = "ImplantCookie"
    cookie_name                = "KLBRSIDdad"
    session_state              = "start"
    session_persistence_period = 3100
  }
  default_forward_rule {
    backend_server_group_id = "49793f55-e5f2-4329-949d-6ea6170ece54"
  }
}
`
