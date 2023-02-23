package ksyun

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKcsSecGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKcsSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKcsSecGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKcsSecurityGroupExists("ksyun_redis_sec_group.test"),
				),
			},
		},
	})
}

func testAccCheckKcsSecurityGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("kcs security group is create failure")
		}
		return nil
	}
}

func testAccCheckKcsSecurityGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*KsyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_redis_sec_group" {
			instanceCheck := make(map[string]interface{})
			instanceCheck["SecurityGroupId"] = rs.Primary.ID
			ptr, err := client.kcsv1conn.DescribeSecurityGroup(&instanceCheck)
			// Verify the error is what we want
			if err != nil {
				return err
			}
			if ptr != nil {
				if (*ptr)["Data"] != nil {
					return errors.New("delete redis sec group failure")
				}
			}
		}
	}

	return nil
}

const testAccKcsSecGroupConfig = `
resource "ksyun_redis_sec_group" "test" {
  name 			= "testTerraform777"
  description 	= "testTerraform777"
}
`
