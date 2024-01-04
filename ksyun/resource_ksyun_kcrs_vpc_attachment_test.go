package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunKcrsVpcAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_kcrs_vpc_attachment.foo",

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKsyunKcrsVpcAttachmentConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_kcrs_vpc_attachment.foo"),
				),
			},
			{
				Config: testAccKsyunKcrsVpcAttachmentUpdateConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_kcrs_vpc_attachment.foo"),
				),
			},
		},
	})
}

const testAccKsyunKcrsVpcAttachmentConfig = `
resource "ksyun_kcrs_vpc_attachment" "foo" {
	instance_id = "b061ebeb-106d-40b9-88ea-7cad7e0c08e5"
	vpc_id = "87e64f91-08eb-405c-8b22-75e075f89aca"
	reserve_subnet_id = "d12b6191-1c4f-433e-b760-419ef79673a3"
	enable_vpc_domain_dns = true
}
`
const testAccKsyunKcrsVpcAttachmentUpdateConfig = `
resource "ksyun_kcrs_vpc_attachment" "foo" {
	instance_id = "b061ebeb-106d-40b9-88ea-7cad7e0c08e5"
	vpc_id = "87e64f91-08eb-405c-8b22-75e075f89aca"
	reserve_subnet_id = "d12b6191-1c4f-433e-b760-419ef79673a3"
	enable_vpc_domain_dns = false
}
`
