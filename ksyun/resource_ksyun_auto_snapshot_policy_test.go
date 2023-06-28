package ksyun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	_ "github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKsyunAutoSnapshot_basic(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_auto_snapshot_policy.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSnapshotDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_auto_snapshot_policy.foo"),
				),
			},
			// to test terraform when its configuration changes
			{
				Config: testAccSnapshotUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_auto_snapshot_policy.foo"),
				),
			},
			{
				Config: testAccSnapshotEmptyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_auto_snapshot_policy.foo"),
				),
			},
		},
	})
}

func testAccCheckSnapshotDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_auto_snapshot_policy" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		snapshotSrv := AutoSnapshotSrv{
			client: client,
		}
		asp := make(map[string]interface{})
		asp["AutoSnapshotPolicyId.0"] = rs.Primary.ID
		sdkResponse, err := snapshotSrv.querySnapshotPolicyByID(asp)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if len(sdkResponse) == 0 {
			continue
		} else {
			return fmt.Errorf("auto snapshot policy still exist")
		}
	}

	return nil
}

const testAccSnapshotConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}


resource "ksyun_auto_snapshot_policy" "foo" {
  name   = "tf_combine_on_hcl"
  auto_snapshot_date = [1,3,4,5]
  auto_snapshot_time = [1,3,4,5,9,22]
}
`

const testAccSnapshotUpdateConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}


resource "ksyun_auto_snapshot_policy" "foo" {
  name   = "tf_combine_on_hcl"
  auto_snapshot_date = [1,3,4,5]
  auto_snapshot_time = [1,3,4,5,9,22]
  retention_time = 3
}
`

// test retention_time: 3 -> 5
// test auto_snapshot_date: [1,3,5] -> [1,4,2], []
// test auto_snapshot_time: [1,3] -> [1,2], []
// retention_time: 5 -> null // persistent save
const testAccSnapshotEmptyConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}


resource "ksyun_auto_snapshot_policy" "foo" {
  name   = "tf_combine_on_hcl_test"
  auto_snapshot_date = [1,3,5]
  auto_snapshot_time = [1,5,6]
  # retention_time = 5
}

`
