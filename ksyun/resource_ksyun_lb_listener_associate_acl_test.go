package ksyun

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pkg/errors"
)

func TestAccKsyunListenerAssociateAcl_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_listener_associate_acl.default",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckListenerAssociateAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlbListenerAssociateAclConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerAssociateAclExists("ksyun_lb_listener_associate_acl.default", &val),
					testAccCheckListenerAssociateAclAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckListenerAssociateAclExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Listener id is empty")
		}
		ids := strings.Split(rs.Primary.ID, ":")
		if len(ids) != 2 {
			return fmt.Errorf("id is error:%v", rs.Primary.ID)
		}
		client := testAccProvider.Meta().(*KsyunClient)
		listener := make(map[string]interface{})
		listener["ListenerId.1"] = ids[0]
		ptr, err := client.slbconn.DescribeListeners(&listener)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["ListenerSet"].([]interface{})
			if len(l) == 0 {
				return errors.New("no listener found")
			}
			/*		listenerItem, ok := l[0].(map[string]interface{})
					if !ok {
						return errors.New("no listener found")
					}
					for key, value := range listenerItem {
						if key == "LoadBalancerAclId" && value == ids[1] {
							val=&map[string]interface{}{
								"LoadBalancerAclId":value,
								"ListenerId":ids[0],
							}
							return nil
						}
					}
			*/
		}
		*val = *ptr
		return nil
	}
}
func testAccCheckListenerAssociateAclAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["ListenerSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("Listener id is empty")
			}
		}
		return nil
	}
}
func testAccCheckListenerAssociateAclDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_lb_listener_associate_acl_associate_acl" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		listener := make(map[string]interface{})
		ids := strings.Split(rs.Primary.ID, ":")
		if len(ids) != 2 {
			return fmt.Errorf("id is error:%v", rs.Primary.ID)
		}
		listener["ListenerId.1"] = ids[0]
		ptr, err := client.slbconn.DescribeListeners(&listener)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["ListenerSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				for _, listener := range l {
					listenerItem, ok := listener.(map[string]interface{})
					if !ok {
						return errors.New("no listener found")
					}
					for key, value := range listenerItem {
						if key == "LoadBalancerAclId" && value == ids[1] {
							return fmt.Errorf("Listener and acl still exist")
						}
					}
				}

			}
		}
	}

	return nil
}

const testAccListenerAssociateAclConfig = `
resource "ksyun_vpc" "default" {
	vpc_name        = "tf-acc-vpc"
	cidr_block = "192.168.0.0/16"
}
resource "ksyun_lb" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  load_balancer_name = "tf-xun-2"
  type = "public"
  subnet_id = ""
  load_balancer_state = "stop"
  private_ip_address = ""
}
# Create Load Balancer Listener with tcp protocol
resource "ksyun_lb_listener" "default" {
  listener_name = "tf-xun-2"
  listener_port = "8080"
  listener_protocol = "HTTP"
  listener_state = "stop"
  load_balancer_id = "${ksyun_lb.default.id}"
  method = "RoundRobin"
  certificate_id = ""
  session {
    session_state = "stop"
    session_persistence_period = 100
    cookie_type = "RewriteCookie"
    cookie_name = "cookiexunqq"
  }
}
resource "ksyun_lb_acl" "default" {
  load_balancer_acl_name = "tf-xun-2"
}
resource "ksyun_lb_listener_associate_acl" "foo" {
  listener_id = "${ksyun_lb_listener.default.id}"
  load_balancer_acl_id = "${ksyun_lb_acl.default.id}"
}
`

const testAccAlbListenerAssociateAclConfig = `
resource "ksyun_lb_listener_associate_acl" "default" {
  listener_id          = "09197c73-88ce-4070-9ec2-2da476742926"
  load_balancer_acl_id = "6de65767-93f7-4b4c-8601-2f7ad28a184e"
  lb_type = "Alb"
}

`
