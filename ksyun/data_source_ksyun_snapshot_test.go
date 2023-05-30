package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunSnapshotDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSnapshotConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_snapshot.foo"),
				),
			},
		},
	})
}

const testAccDataSourceSnapshotConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}

data "ksyun_snapshot" "foo" {
	name = "tf_combine_test"
	auto_snapshot_policy_ids = ["858469026661474304"]
    output_file = "output_result_snapshot"
}

output "ksyun_snapshot" {
	value = data.ksyun_snapshot.foo
}

output "ksyun_snapshots_total_count" {
	value = data.ksyun_snapshot.foo.total_count
}
`
