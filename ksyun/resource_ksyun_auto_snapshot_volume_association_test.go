// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
			// // to test terraform when its configuration changes
			// {
			// 	Config: testAccUpdateSnapshotConfig,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckIDExists("ksyun_snapshot.foo"),
			// 	),
			// },
		},
	})
}

const testAccAutoSnapshotVolumeAssociationConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}


resource "ksyun_auto_snapshot_volume_association" "foo" {
  attach_volume_id = "a94bbaac-0b83-4610-9040-cdbc15b061ab"
  auto_snapshot_policy_id = "860576274707722240"
}
`
