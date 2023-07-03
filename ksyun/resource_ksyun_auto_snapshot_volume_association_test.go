package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestResourceKsyunAutoSnapshotAssociation_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_auto_snapshot_volume_association.foo",
		Providers:     testAccProviders,
		// CheckDestroy:  testAccCheckSnapshotDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccAutoSnapshotVolumeAssociationConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_auto_snapshot_volume_association.foo"),
				),
			},
			// to test terraform when its configuration changes
			{
				Config: testAccAutoSnapshotVolumeAssociationUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_auto_snapshot_volume_association.foo"),
				),
			},
		},
	})
}

const testAccAutoSnapshotVolumeAssociationConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}


resource "ksyun_auto_snapshot_volume_association" "foo" {
  attach_volume_id = "6a636034-64f9-4fb6-8336-485ecf47fc65"
  auto_snapshot_policy_id = "860576274707722240"
}
`
const testAccAutoSnapshotVolumeAssociationUpdateConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}


resource "ksyun_auto_snapshot_volume_association" "foo" {
  attach_volume_id = "9f02901e-2573-454c-9250-009507a030be"
  auto_snapshot_policy_id = "860576274707722240"
}
`
