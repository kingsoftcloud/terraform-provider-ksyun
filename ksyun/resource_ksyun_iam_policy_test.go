package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunIamPolicy_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIAMPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_iam_policy.policy"),
				),
			},
			{
				Config: testAccIAMPolicyConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_iam_policy.policy"),
				),
			},
		},
	})
}

const testAccIAMPolicyConfig = `
resource "ksyun_iam_policy" "policy" {
  policy_name = "TestPolicy1"
  policy_document = "{\"Version\": \"2015-11-01\",\"Statement\": [{\"Effect\": \"Allow\",\"Action\": [\"iam:ListUsers\"],\"Resource\": [\"*\"]}]}"
}`

const testAccIAMPolicyConfigUpdate = `
resource "ksyun_iam_policy" "policy" {
  policy_name = "TestPolicy1"
  policy_document = "{\"Version\": \"2015-11-01\",\"Statement\": [{\"Effect\": \"Allow\",\"Action\": [\"iam:List*\"],\"Resource\": [\"*\"]}]}"
}`
