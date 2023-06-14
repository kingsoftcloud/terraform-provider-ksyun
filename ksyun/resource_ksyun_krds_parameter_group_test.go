// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
	region =  "cn-qingyangtest-1"
}

resource "ksyun_krds_parameter_group" "dpg_with_parameters" {
  name  = "tf_krdpg_on_hcl_with"
  description    = "acceptance-test"
  engine = "mysql"
  engine_version = "5.5"
  parameters = {
    auto_increment_increment = 10240
    auto_increment_offset = 5
    back_log = 65535
	connect_timeout = 30
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
  parameters = {
	    auto_increment_increment = 10240
		auto_increment_offset = 5
		back_log = 65535
		connect_timeout = 60
		table_open_cache_instances = 1  
		group_concat_max_len = 102
	}
}
`
const testAccKrdsParameterGroupApplyConfig = `

`
