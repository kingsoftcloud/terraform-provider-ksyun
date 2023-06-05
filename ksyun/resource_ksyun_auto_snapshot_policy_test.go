package ksyun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	_ "github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKsyunSnapshot_basic(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_snapshot.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSnapshotDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_snapshot.foo"),
				),
			},
			// to test terraform when its configuration changes
			{
				Config: testAccUpdateSnapshotConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_snapshot.foo"),
				),
			},
		},
	})
}

// func TestAccKsyunSubnet_update(t *testing.T) {
// 	var val map[string]interface{}
//
// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck: func() {
// 			testAccPreCheck(t)
// 		},
//
// 		IDRefreshName: "ksyun_subnet.foo",
// 		Providers:     testAccProviders,
// 		CheckDestroy:  testAccCheckSubnetDestroy,
//
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccSubnetConfig,
//
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckSubnetExists("ksyun_subnet.foo", &val),
// 					testAccCheckSubnetAttributes(&val),
// 				),
// 			},
// 			{
// 				Config: testAccSubnetUpdateConfig,
//
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckSubnetExists("ksyun_subnet.foo", &val),
// 					testAccCheckSubnetAttributes(&val),
// 				),
// 			},
// 		},
// 	})
// }
//
// func testAccCheckSubnetExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[n]
//
// 		if !ok {
// 			return fmt.Errorf("not found: %s", n)
// 		}
//
// 		if rs.Primary.ID == "" {
// 			return fmt.Errorf("Snapshot id is empty")
// 		}
//
// 		client := testAccProvider.Meta().(*KsyunClient)
// 		subnet := make(map[string]interface{})
// 		subnet["SubnetId.1"] = rs.Primary.ID
// 		ptr, err := client.vpcconn.DescribeSubnets(&subnet)
//
// 		if err != nil {
// 			return err
// 		}
// 		if ptr != nil {
// 			l := (*ptr)["SubnetSet"].([]interface{})
// 			if len(l) == 0 {
// 				return err
// 			}
// 		}
//
// 		*val = *ptr
// 		return nil
// 	}
// }
//
// func testAccCheckSubnetAttributes(val *map[string]interface{}) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		if val != nil {
// 			l := (*val)["SubnetSet"].([]interface{})
// 			if len(l) == 0 {
// 				return fmt.Errorf("subnet id is empty")
// 			}
// 		}
// 		return nil
// 	}
// }
//
func testAccCheckSnapshotDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_snapshot" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		snapshotSrv := AutoSnapshotSrv{
			client: client,
		}
		subnet := make(map[string]interface{})
		subnet["AutoSnapshotPolicyId.1"] = rs.Primary.ID
		ptr, err := snapshotSrv.querySnapshotPolicyByID(subnet)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			policySetIf := ptr["AutoSnapshotPolicySet"]
			if policySetIf == nil {
				continue
			}
			policySet := policySetIf.([]interface{})
			if len(policySet) == 0 {
				continue
			} else {
				return fmt.Errorf("subnet still exist")
			}
		}
	}

	return nil
}

const testAccSnapshotConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}


resource "ksyun_snapshot" "foo" {
  name   = "tf_combine_on_hcl"
  auto_snapshot_date = [1,3,4,5]
  auto_snapshot_time = [1,3,4,5,9,22]
}
`

const testAccUpdateSnapshotConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}


resource "ksyun_snapshot" "foo" {
  name   = "tf_combine_on_hcl"
  auto_snapshot_date = [1,3,4,5]
  auto_snapshot_time = [1,3,4,5,9,22]
  retention_time = 3
}
`
