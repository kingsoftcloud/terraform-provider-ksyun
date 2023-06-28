package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunAutoSnapshotVolumeAssociationDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// {
			// 	Config: testAccDataKrdsParameterGroupConfig,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckIDExists("data.ksyun_auto_snapshot_policy_volume_association.foo"),
			// 	),
			// },
			{
				Config: testAccDataASP2VAssociationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_auto_snapshot_policy_volume_association.foo"),
				),
			},
		},
	})
}

const testAccDataASP2VAssociationConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}

data "ksyun_auto_snapshot_policy_volume_association" "foo" {
	output_file = "output_result_volume_id"
	attach_volume_id = "9e12a289-1653-4524-8c2b-aa5e3fe557c0"
}

data "ksyun_auto_snapshot_policy_volume_association" "foo1" {
	output_file = "output_result_null"
}
data "ksyun_auto_snapshot_policy_volume_association" "foo2" {
	output_file = "output_result_policy_id"
	auto_snapshot_policy_id = "860576274707722240"
}
`
