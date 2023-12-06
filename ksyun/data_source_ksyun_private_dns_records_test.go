package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunPdnsRecordsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPdnsrecordsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_private_dns_records.foo"),
				),
			},
		},
	})
}

const testAccDataPdnsrecordsConfig = `

data "ksyun_private_dns_records" "foo" {
  output_file = "pdns_records_output_result"
  zone_id = "a5ae6bf0-0ff4-472f-af9e-158afab99915"
  region_name = ["cn-beijing-6"]
}
`
