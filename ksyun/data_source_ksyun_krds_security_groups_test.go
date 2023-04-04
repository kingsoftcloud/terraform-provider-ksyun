package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunKrdsSecurityGroupDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKrdsSecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_krds_security_groups.search-krds-sec-group"),
				),
			},
		},
	})
}

const testAccDataKrdsSecurityGroupConfig = `
resource "ksyun_krds_security_group" "krds_sec_group_14" {
  security_group_name = "terraform_security_group_14"
  security_group_description = "terraform-security-group-14"
  security_group_rule{
    security_group_rule_protocol = "182.133.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.134.0.0/16"
    security_group_rule_name = "wtf"
  }

}

data "ksyun_krds_security_groups" "search-krds-sec-group"{
  output_file = "output_file"
  security_group_id = "${ksyun_krds_security_group.krds_sec_group_14.id}"
}
`
