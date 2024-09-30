package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunIamProject_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIAMProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_iam_project.project"),
				),
			},
		},
	})
}

const testAccIAMProjectConfig = `
resource "ksyun_iam_project" "project" {
  project_name = "ProjectNameTest"
  project_desc = "ProjectDescTest"
}`
