package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunKmrClusters_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKmrClustersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kmr_clusters.default"),
					resource.TestCheckResourceAttr("data.ksyun_kmr_clusters.default", "clusters.#", "10"),         // 检查clusters数量
					resource.TestCheckResourceAttrSet("data.ksyun_kmr_clusters.default", "clusters.0.cluster_id"), // 检查有数据
				),
			},
		},
	})
}

const testAccDataKmrClustersConfig = `
provider "ksyun" {
  access_key = "your ak"
  secret_key = "your sk"
  region     = "cn-beijing-6"
  endpoint   = "kmr.api.ksyun.com"
}

data "ksyun_kmr_clusters" "default" {
  marker      = "limit=10&offset=0"
  output_file = "output_result"
}

output "total" {
  value = data.ksyun_kmr_clusters.default.total
}
`
