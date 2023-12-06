package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunPrivateDnsVpcAttachment_basic(t *testing.T) {
	// var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_private_dns_zone_vpc_attachment.example",
		Providers:     testAccProviders,
		// CheckDestroy:  testAccCheckPrivateDnsDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsVpcAttachmentConfig,
				// Check: resource.ComposeTestCheckFunc(
				// 	testAccCheckPrivateDnsExists("ksyun_private_dns_VpcAttachment.foo", &val),
				// 	testAccCheckPrivateDnsAttributes(&val),
				// ),

			},
			{
				ResourceName:      "ksyun_private_dns_zone_vpc_attachment.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPrivateDnsVpcAttachmentConfig = `
provider "ksyun" {
  region = "cn-beijing-6"
}

resource "ksyun_private_dns_zone_vpc_attachment" "example" {
	zone_id = "9e680c2b-bf22-4511-820a-1e2e968667b6"
	vpc_set {
		region_name = "cn-beijing-6"
		vpc_id = "fc592cff-bb85-40bb-9ac3-1218b2ac3318"
    }
}
`
