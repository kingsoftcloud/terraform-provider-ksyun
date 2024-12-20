package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunIamRelationPolicy_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKsyunIamRelationPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_iam_relation_policy.user"),
				),
			},
		},
	})
}

const testAccKsyunIamRelationPolicyConfig = `
resource "ksyun_iam_relation_policy" "user" {
  name = "username01"
  policy_name = "IAMReadOnlyAccess"
  relation_type = 1
  policy_type = "system"
}`
