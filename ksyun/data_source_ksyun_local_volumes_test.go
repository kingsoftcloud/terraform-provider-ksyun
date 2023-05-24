package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunInstanceLocalVolumesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataInstanceLocalVolumesConfig,
				Check: resource.ComposeTestCheckFunc(
					//func(s *terraform.State) error {
					//	n := "data.ksyun_instance_local_volumes.foo"
					//	rs, ok := s.RootModule().Resources[n]
					//	if !ok {
					//		return fmt.Errorf(" Can't find resource or data source: %s ", n)
					//	}
					//	//fmt.Println(s, rs, "123")
					//	if rs.Primary.ID == "" {
					//		return fmt.Errorf("ID is not be set")
					//	}
					//	return nil
					//},
					testAccCheckIDExists("data.ksyun_instance_local_volumes.foo"),
				),
			},
		},
	})
}

const testAccDataInstanceLocalVolumesConfig = `
data "ksyun_instance_local_volumes" "foo" {
  output_file = "output_result"
}`
