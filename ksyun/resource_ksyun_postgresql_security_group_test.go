package ksyun

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunPostgresqlSecrityGroup_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_postgresql_security_group.postgresql_sec_group_239",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckPostgresqlSecrityGroupDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlSecrityGroupConfig,

				Check: resource.ComposeTestCheckFunc(
					testCheckPostgresqlSecrityGroupExists("ksyun_postgresql_security_group.postgresql_sec_group_239", &val),
				),
			},
		},
	})
}

func testCheckPostgresqlSecrityGroupExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found : %s", n)
		}
		if res.Primary.ID == "" {
			return fmt.Errorf("instance is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		req := map[string]interface{}{
			"SecurityGroupId": res.Primary.ID,
		}
		resp, err := client.postgresqlconn.DescribeSecurityGroup(&req)
		if err != nil {
			return err
		}
		*val = *resp
		return nil
	}
}

func testAccCheckPostgresqlSecrityGroupDestroy(s *terraform.State) error {
	for _, res := range s.RootModule().Resources {
		if res.Type != "ksyun_postgresql_security_group" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		req := map[string]interface{}{
			"SecurityGroupId": res.Primary.ID,
		}
		_, err := client.postgresqlconn.DeleteSecurityGroup(&req)
		if err != nil {
			if err.(awserr.Error).Code() == "NOT_FOUND" {
				return nil
			}
			return err
		}
	}

	return nil
}

const testAccPostgresqlSecrityGroupConfig = `


resource "ksyun_postgresql_security_group" "postgresql_sec_group_239" {
  security_group_name = "terraform_security_group_239"
  security_group_description = "terraform-security-group-239"
  security_group_rule{
    security_group_rule_protocol = "182.133.0.1/16"
    security_group_rule_name = "asdf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.134.0.0/16"
    security_group_rule_name = "asdf2"
  }
}



`
