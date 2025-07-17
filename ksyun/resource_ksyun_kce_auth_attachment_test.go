package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const testAccKceAuthAttachment = `
resource "ksyun_kce_auth_attachment" "auth" {
  sub_user_id = "38435"
  permissions {

  cluster_id = "4cf5b24b-de39-4f55-b0ce-fd7b28cb964c"
  cluster_role = "kce:dev"
  namespace = ""
}
}
`

func TestAccKsyunKceAuthAttachment_basic(t *testing.T) {
	//var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_kce_auth_attachment.auth",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKceAuthAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				// PlanOnly: true,
				Config: testAccKceAuthAttachment,
				//Check: resource.ComposeTestCheckFunc(
				//	testAccCheckKceAuthAttachmentExists("ksyun_kce_Auth_attachment.foo", &val),
				//	testAccCheckKceAuthAttachmentAttributes(&val),
				//),
			},
			// {
			//	ResourceName:      "ksyun_kce_AuthAttachment.foo",
			//	ImportStateId:     "bedfb5d0-bb8f-40dd-9f1d-8966fd1ace87:3a66ad5a-313b-41c0-8b72-6749e438ea17",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	//ImportStateVerifyIgnore: []string{"lang"},
			// },
		},
	})
}
func testAccCheckKceAuthAttachmentExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}
func testAccCheckKceAuthAttachmentAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}
func testAccCheckKceAuthAttachmentDestroy(s *terraform.State) error {
	return nil
}
