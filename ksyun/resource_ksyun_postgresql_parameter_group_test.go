package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestResourceKsyunPostgresqlParameterGroup_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_postgresql_parameter_group.dpg_with_parameters",
		Providers:     testAccProviders,
		// CheckDestroy:  testAccCheckSnapshotDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlParameterGroupConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_postgresql_parameter_group.dpg_with_parameters"),
				),
			},
			// to test terraform when its configuration changes
			{
				Config: testAccPostgresqlParameterGroupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_postgresql_parameter_group.dpg_with_parameters"),
				),
			},
		},
	})
}

const testAccPostgresqlParameterGroupConfig = `
provider "ksyun" {
	region =  "cn-beijing-6"
}

resource "ksyun_postgresql_parameter_group" "dpg_with_parameters" {
  name  = "tf_krdpg_on_hcl_with"
  description    = "acceptance-test"
  engine = "postgresql"
  engine_version = "15"
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

const testAccPostgresqlParameterGroupUpdateConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}

// test parameters: -> empty
resource "ksyun_postgresql_parameter_group" "dpg_with_parameters" {
  name  = "tf_krdpg_on_hcl_with"
  description    = "acceptance-test"
  engine = "postgresql"
  engine_version = "15"
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
const testAccPostgresqlParameterGroupApplyConfig = `

`
