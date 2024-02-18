package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunKS3DataSourceInfo(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKS3BucketConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_ks3_buckets.default"),
				),
			},
		},
	})
}

const testAccDataKS3BucketConfig = `
data "ksyun_ks3_buckets" "default" {
  name_regex  = "bucket-202402"
  output_file = "bucket_info.txt"
}
`
