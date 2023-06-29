package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunAutoSnapshotPolicyDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSnapshotConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_auto_snapshot_policy.foo"),
				),
			},
		},
	})
}

const testAccDataSourceSnapshotConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}

data "ksyun_auto_snapshot_policy" "foo" {
    output_file = "output_result_snapshot"
}

output "ksyun_auto_snapshot_policy" {
	value = data.ksyun_auto_snapshot_policy.foo
}

output "ksyun_auto_snapshot_policy_total_count" {
	value = data.ksyun_auto_snapshot_policy.foo.total_count
}
`
