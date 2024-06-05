package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const testAccKceWorkerConfig = `
resource "ksyun_kce_cluster_attach_existence" "foo" {
	cluster_id = "45e21f7e-fd87-4c45-9e58-e3e2641b0729"
	instance_id = "d9d852da-9e04-4fc3-a2ae-24950d97a167"
	image_id = "7dc43a49-4d3e-4498-993c-4192847d75bf"

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
		IDRefreshName: "ksyun_kce_cluster_attach_existence.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKceWorkerDestroy,
		Steps: []resource.TestStep{
			{
				// PlanOnly: true,
				Config: testAccKceWorkerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKceWorkerExists("ksyun_kce_cluster_attach_existence.foo", &val),
					testAccCheckKceWorkerAttributes(&val),
				),
			},
			// {
			//	ResourceName:      "ksyun_kce_worker.foo",
			//	ImportStateId:     "bedfb5d0-bb8f-40dd-9f1d-8966fd1ace87:3a66ad5a-313b-41c0-8b72-6749e438ea17",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	//ImportStateVerifyIgnore: []string{"lang"},
			// },
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
