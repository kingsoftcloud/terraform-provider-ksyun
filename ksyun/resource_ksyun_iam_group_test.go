package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunIamGroup_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIAMGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_iam_group.group"),
				),
			},
		},
	})
}

const testAccIAMGroupConfig = `
resource "ksyun_iam_group" "group" {
  group_name = "GroupNameTest"
  description = "desc"
}`
