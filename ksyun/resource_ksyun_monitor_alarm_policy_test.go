package ksyun

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func TestAccKsyunMonitorAlarmPolicy_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorAlarmPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorAlarmPolicyExists("ksyun_monitor_alarm_policy.foo"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "policy_name", "tf-test-policy"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "product_type", "0"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "policy_type", "0"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "resource_bind_type", "3"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "trigger_rules.#", "1"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "trigger_rules.0.period", "5m"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "trigger_rules.0.method", "avg"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "trigger_rules.0.compare", ">"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "trigger_rules.0.trigger_value", "90"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "trigger_rules.0.item_name", "CPU利用率"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "trigger_rules.0.item_key", "cpu.utilizition.total"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "trigger_rules.0.units", "%"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "trigger_rules.0.points", "2"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "trigger_rules.0.interval", "5"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "trigger_rules.0.max_count", "3"),
				),
			},
		},
	})
}

func TestAccKsyunMonitorAlarmPolicy_withUserNotice(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorAlarmPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmPolicyWithUserNoticeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorAlarmPolicyExists("ksyun_monitor_alarm_policy.foo"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "policy_name", "tf-test-policy-notice"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "user_notice.#", "1"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "user_notice.0.contact_way", "3"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "user_notice.0.contact_flag", "1"),
				),
			},
		},
	})
}

func TestAccKsyunMonitorAlarmPolicy_withUrlNotice(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorAlarmPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmPolicyWithUrlNoticeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorAlarmPolicyExists("ksyun_monitor_alarm_policy.foo"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "policy_name", "tf-test-policy-url"),
					resource.TestCheckResourceAttr("ksyun_monitor_alarm_policy.foo", "url_notice.#", "1"),
				),
			},
		},
	})
}

func testAccCheckMonitorAlarmPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("monitor alarm policy create failure")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		readReq := make(map[string]interface{})
		readReq["PolicyId"] = rs.Primary.ID

		logger.Debug(logger.ReqFormat, "DescribeAlarmPolicy", readReq)
		_, err := client.monitorv4conn.DescribeAlarmPolicy(&readReq)
		if err != nil {
			return fmt.Errorf("error on reading monitor alarm policy %q, %s", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccCheckMonitorAlarmPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*KsyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_monitor_alarm_policy" {
			policyCheck := make(map[string]interface{})
			policyCheck["PolicyId"] = rs.Primary.ID
			_, err := client.monitorv4conn.DescribeAlarmPolicy(&policyCheck)

			if err != nil {
				if strings.Contains(err.Error(), "PolicyNotFound") || strings.Contains(err.Error(), "not found") {
					return nil
				} else {
					return fmt.Errorf("monitor alarm policy delete failure: %s", err)
				}
			} else {
				return fmt.Errorf("monitor alarm policy still exists")
			}
		}
	}

	return nil
}

const testAccMonitorAlarmPolicyConfig = `
resource "ksyun_monitor_alarm_policy" "foo" {
  policy_name        = "tf-test-policy"
  product_type       = 0
  policy_type       = 0
  resource_bind_type = 3

  trigger_rules {
    compare       = ">"
	effect_bt     = "00:00"
	effect_et     = "23:59"
    interval      = 5
	item_key      = "cpu.utilizition.total"
    item_name     = "CPU利用率"
	max_count     = 3
    method        = "avg"
    period        = "5m"
    points        = 2
    trigger_value = "90"
    units         = "%"
  }
	
  instance_ids = []

  user_notice {
    contact_way  = 2
    contact_flag = 2
    contact_id   = 4423
  }

  url_notice = ["https://example.com/webhook"]
}
`

const testAccMonitorAlarmPolicyWithUserNoticeConfig = `
resource "ksyun_monitor_alarm_policy" "foo" {
  policy_name        = "tf-test-policy-notice"
  product_type       = 0
  policy_type        = 0
  resource_bind_type = 3

  trigger_rules {
    compare       = ">"
	effect_bt     = "00:00"
	effect_et     = "23:59"
    interval      = 5
	item_key      = "cpu.utilizition.total"
    item_name     = "CPU利用率"
	max_count     = 3
    method        = "avg"
    period        = "5m"
    points        = 2
    trigger_value = "90"
    units         = "%"
  }

  instance_ids = []

  user_notice {
    contact_way  = 2
    contact_flag = 2
    contact_id   = 4423
  }

  url_notice = ["https://example.com/webhook"]
}
`

const testAccMonitorAlarmPolicyWithUrlNoticeConfig = `
resource "ksyun_monitor_alarm_policy" "foo" {
  policy_name        = "tf-test-policy-url"
  product_type       = 0
  policy_type        = 0
  resource_bind_type = 3
  
  instance_ids = []

  trigger_rules {
    compare       = ">"
	effect_bt     = "00:00"
	effect_et     = "23:59"
    interval      = 5
	item_key      = "cpu.utilizition.total"
    item_name     = "CPU利用率"
	max_count     = 3
    method        = "avg"
    period        = "5m"
    points        = 2
    trigger_value = "90"
    units         = "%"
  }

  url_notice = ["https://example.com/webhook"]
}
`
