package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const testAccKceClusterAttachmentConfigWithNewInstance = `
resource "ksyun_kce_cluster_attachment" "foo" {
	cluster_id = "dec547af-a10d-4f21-82b4-89ff5642c55a"
	
  worker_config {
    image_id      = "fbafd8cd-b570-47c4-a3db-ff9702108f17"
    instance_type = "S6.4B"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = "c771027a-fafd-4b3b-a6b9-daeab9d0c13a"
    security_group_id = ["59a87036-dc27-41cf-98ab-24a387501195"]
    charge_type       = "Daily"

  }
	advanced_setting {
	  container_runtime = "containerd"
	  pre_user_script   = "def"
      label {
		key  = "tf_assembly_kce"
		value = "advanced_setting"
      }
	  taints {
		key    = "key1"
		value  = "value1"
		effect = "NoSchedule"

      }
	  extra_arg = ["--feature-gates=EphemeralContainers=true", "--allow-privileged=true"]
	}
}
`

func TestAccKsyunKceClusterAttachment_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_kce_cluster_attachment.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKceClusterAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				// PlanOnly: true,
				Config: testAccKceClusterAttachmentConfigWithNewInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKceClusterAttachmentExists("ksyun_kce_cluster_attachment.foo", &val),
					testAccCheckKceClusterAttachmentAttributes(&val),
				),
			},
			// {
			//	ResourceName:      "ksyun_kce_ClusterAttachment.foo",
			//	ImportStateId:     "bedfb5d0-bb8f-40dd-9f1d-8966fd1ace87:3a66ad5a-313b-41c0-8b72-6749e438ea17",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	//ImportStateVerifyIgnore: []string{"lang"},
			// },
		},
	})
}
func testAccCheckKceClusterAttachmentExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}
func testAccCheckKceClusterAttachmentAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}
func testAccCheckKceClusterAttachmentDestroy(s *terraform.State) error {
	return nil
}
