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
provider "ksyun" {
  #指定KS3服务的访问域名
  endpoint = "ks3-cn-beijing.ksyuncs.com"
}

data "ksyun_ks3_buckets" "default" {
  #匹配全部包涵该字符串的bucket
  name_regex  = "bucket-202402"
  #输出文件路径
  output_file = "bucket_info.txt"
}
`
