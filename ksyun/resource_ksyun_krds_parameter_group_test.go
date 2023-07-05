package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestResourceKsyunKrdsParameterGroup_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_krds_parameter_group.dpg_with_parameters",
		Providers:     testAccProviders,
		// CheckDestroy:  testAccCheckSnapshotDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccKrdsParameterGroupConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_krds_parameter_group.dpg_with_parameters"),
				),
			},
			// to test terraform when its configuration changes
			{
				Config: testAccKrdsParameterGroupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_krds_parameter_group.dpg_with_parameters"),
				),
			},
		},
	})
}

const testAccKrdsParameterGroupConfig = `
provider "ksyun" {
	region =  "cn-beijing-6"
}

resource "ksyun_krds_parameter_group" "dpg_with_parameters" {
  name  = "tf_krdpg_on_hcl_with"
  description    = "acceptance-test"
  engine = "mysql"
  engine_version = "5.7"
parameters {
	    	name = "auto_increment_increment"
	    	value = "8"
		}
		parameters {
			name = "binlog_format"
			value = "ROW"
		}
		parameters {
			name = "delayed_insert_limit"
			value = "108"
		}
		parameters {
			name = "auto_increment_offset"
			value= "2"
		}
		parameters {
			name = "innodb_open_files"
			value= "1000"
		}
}
`

const testAccKrdsParameterGroupUpdateConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}

// test parameters: -> empty
resource "ksyun_krds_parameter_group" "dpg_with_parameters" {
  name  = "tf_krdpg_on_hcl_with"
  description    = "acceptance-test"
  engine = "mysql"
  engine_version = "5.7"
parameters {
	    	name = "auto_increment_increment"
	    	value = "8"
		}
		parameters {
			name = "binlog_format"
			value = "ROW"
		}
		parameters {
			name = "delayed_insert_limit"
			value = "108"
		}
		parameters {
			name = "auto_increment_offset"
			value= "2"
		}
		parameters {
			name = "back_log"
			value= "65535"
		}
		parameters {
			name = "innodb_open_files"
			value= "1024"
		}
}
`
const testAccKrdsParameterGroupApplyConfig = `

`
