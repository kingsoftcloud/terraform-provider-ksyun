package ksyun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunTagV2Attachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckVolumeAttachDestory,
		Steps: []resource.TestStep{
			{
				Config: tagv2Config,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccCheckTagv2AttachExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("attach id is empty")
		}
		return nil
	}
}

const tagv2Config = `
variable "bucket_names" {
  default = {
    "test_tag_key_2" = "1c415524-b846-4c40-a8dc-5159055dd457"
    "test_tag_key"   = "1c415524-b846-4c40-a8dc-5159055dd457"
  }
}

resource "ksyun_tag_v2_attachment" "kec_tag1" {
  for_each      = var.bucket_names
  key           = each.key
  value         = "test_tag_value"
  resource_type = "bws"
  resource_id   = each.value

  # depends_on = [ksyun_tag_v2.tagv2]
}



`
