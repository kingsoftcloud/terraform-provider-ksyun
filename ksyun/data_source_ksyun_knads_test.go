package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunKnadsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKnadsConfig,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["data.ksyun_knads.foo"]

						if !ok {
							return fmt.Errorf(" Can't find resource or data source: %s ", "data.ksyun_knads.foo")
						}

						if rs.Primary.ID == "" {
							return fmt.Errorf("ID is not be set")
						}
						return nil
					}),
			},
		},
	})
}

const testAccDataKnadsConfig = `

data "ksyun_knads" "foo" {
 ids = ["knad3f694ff7-83fe-34a4-818c-45ad7da65ef6","knadb4fe2700-8565-3266-a51a-c37a2987f73b"]
project_id = ["0"]
  /*service_id= "KEAD_100G"
  band=30
  max_band=30
  ip_count= 10
  bill_type= 1
  idc_band= 100*/
}
`
const testAccDataKnadsUpdateConfig = `

data "ksyun_knads" "foo" {
 
  service_id= "KEAD_50G"
  band=30
  max_band=30
  ip_count= 10
  bill_type= 1
  idc_band= 100
}
`
