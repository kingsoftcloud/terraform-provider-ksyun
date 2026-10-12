package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunKesClusters_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKesClustersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kes_clusters.default"),
					resource.TestCheckResourceAttrSet("data.ksyun_kes_clusters.default", "total"),
				),
			},
		},
	})
}

const testAccDataKesClustersConfig = `
provider "ksyun" {
  access_key = "your ak"
  secret_key = "your sk"
  region     = "cn-beijing-6"
  endpoint   = "kes.api.ksyun.com"
}

data "ksyun_kes_clusters" "default" {
  marker      = "limit=10&offset=0"
  output_file = "output_result"
}

output "total" {
  value = data.ksyun_kes_clusters.default.total
}
`
