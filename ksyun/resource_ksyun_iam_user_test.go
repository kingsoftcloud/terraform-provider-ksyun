package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunIamUser_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIAMUserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_iam_user.user"),
				),
			},
		},
	})
}

const testAccIAMUserConfig = `
resource "ksyun_iam_user" "user" {
  user_name = "username01"
  real_name = "realname01"
  phone = "13800000000"
  email = "test@ksyun.com"
  remark = "remark"
  password = "password"
  password_reset_required = 0
  open_login_protection = 1
  open_security_protection = 1
  view_all_project = 0
}`
