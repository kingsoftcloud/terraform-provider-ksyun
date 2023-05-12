package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

const testAccKceWorkerConfig = `
resource "ksyun_kce_worker" "foo" {
	cluster_id = "bdd98fcf-b3d3-4296-817a-f54b04f9ea9e"
	image_id = "234567"
	instance_id = "7496b61d-8d50-4c5a-86b3-4953cf63f3fd"
}
`

func TestAccKsyunKceWorker_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_kce_worker.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKceWorkerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKceWorkerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKceWorkerExists("ksyun_kce_worker.foo", &val),
					testAccCheckKceWorkerAttributes(&val),
				),
			},
		},
	})
}
func testAccCheckKceWorkerExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}
func testAccCheckKceWorkerAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}
func testAccCheckKceWorkerDestroy(s *terraform.State) error {
	return nil
}
