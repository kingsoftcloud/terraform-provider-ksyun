package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunKesCluster_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKesClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_kes_cluster.default"),
					resource.TestCheckResourceAttrSet("data.ksyun_kes_cluster.default", "clusters.0.cluster_id"),
					resource.TestCheckResourceAttrSet("data.ksyun_kes_cluster.default", "clusters.0.cluster_name"),
				),
			},
		},
	})
}

const testAccDataKesClusterConfig = `
provider "ksyun" {
  access_key = "your ak"
  secret_key = "your sk"
  region     = "cn-beijing-6"
  endpoint   = "kes.api.ksyun.com"
}

data "ksyun_kes_cluster" "default" {
  cluster_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  output_file = "output_result"
}

output "cluster_name" {
  value = data.ksyun_kes_cluster.default.clusters.0.cluster_name
}
`
