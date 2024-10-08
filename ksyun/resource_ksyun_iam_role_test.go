package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunIamRole_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIAMRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_iam_role.role"),
				),
			},
		},
	})
}

const testAccIAMRoleConfig = `
resource "ksyun_iam_role" "role" {
  role_name = "role_name_test"
  trust_accounts = "2000096256"
  description = "desc"
}`
