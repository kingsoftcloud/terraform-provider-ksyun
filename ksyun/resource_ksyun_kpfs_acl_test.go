package ksyun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestKpfsAcl_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_kpfs_acl.kpfs-acl-1",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					testAclUpdate("ksyun_kpfs_acl.kpfs-acl-1", &val),
					// testAccCheckLbAttributes(&val),
				),
			},
			// {
			// 	Config: testConfig,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckLbExists("ksyun_lb.foo", &val),
			// 		testAccCheckLbAttributes(&val),
			// 	),
			// },
		},
	})
}

const testConfig = `
resource "ksyun_kpfs_acl" "kpfs-acl-1"{
	epc_id = "c6c683f8-5bb4-4747-8516-9a61f01c4bce"
	kpfs_acl_id = "c6c683f8-5bb4-4747-8516-9a61f01c4bce"
  }
`

func testAclUpdate(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		client := testAccProvider.Meta().(*KsyunClient)
		req := make(map[string]interface{})
		req["PosixAclId"] = rs.Primary.Attributes["kpfs_acl_id"]
		ptr, err := client.kpfsconn.DescribePerformanceOnePosixAclList(&req)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["DescribePerformanceOnePosixAclList"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}
