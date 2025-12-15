package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunKlogProjectsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKlogProjectsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_klog_projects.foo"),
				),
			},
		},
	})
}

const testAccDataKlogProjectsConfig = `

data "ksyun_klog_projects" "foo" {
  project_name=""
  output_file = "output_result"
}
`
