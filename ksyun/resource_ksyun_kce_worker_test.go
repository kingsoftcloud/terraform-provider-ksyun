package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

const testAccKceWorkerConfig = `
resource "ksyun_kce_worker" "foo" {
	cluster_id = "5e2d7073-a948-42d1-81ec-13d679edfd10"
	image_id = "7dc43a49-4d3e-4498-993c-4192847d75bf"
	instance_id = "321521e8-4885-426c-baa1-2fc9a4eced76"

	instance_password = "Test1234$"

	data_disk {
		auto_format_and_mount = true
		file_system = "ext4"
		mount_target = "/data"
	}
	container_runtime = "docker"
	docker_path = "/data/docker_new"
	user_script = "abc"
	pre_user_script = "def"
	schedulable = true
	//label {
	//	key = "key1"
	//	value = "value1"
	//}
	//label {
	//	key = "key2"
	//	value = "value2"
	//}
	container_log_max_size = 200
	container_log_max_files = 20
	extra_arg = ["abc=def", "hig=klm"]
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
				//PlanOnly: true,
				Config: testAccKceWorkerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKceWorkerExists("ksyun_kce_worker.foo", &val),
					testAccCheckKceWorkerAttributes(&val),
				),
			},
			//{
			//	ResourceName:      "ksyun_kce_worker.foo",
			//	ImportStateId:     "bedfb5d0-bb8f-40dd-9f1d-8966fd1ace87:3a66ad5a-313b-41c0-8b72-6749e438ea17",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	//ImportStateVerifyIgnore: []string{"lang"},
			//},
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
