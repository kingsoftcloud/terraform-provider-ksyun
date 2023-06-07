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

		IDRefreshName: "ksyun_krds_parameter_group.foo",
		Providers:     testAccProviders,
		// CheckDestroy:  testAccCheckSnapshotDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccKrdsParameterGroupConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_krds_parameter_group.foo"),
				),
			},
			// // to test terraform when its configuration changes
			{
				Config: testAccKrdsParameterGroupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_krds_parameter_group.foo"),
				),
			},
		},
	})
}

const testAccKrdsParameterGroupConfig = `
provider "ksyun" {
	region =  "cn-beijing-6"
}

resource "ksyun_krds_parameter_group" "foo" {
	name = "tf_test_on_hcl"
	engine = "mysql"
	engine_version = "5.7"
	description = "terraform created"
	parameters = {
		connect_timeout = 30
		log_slow_admin_statements = "OFF"
		log_bin_trust_function_creators = "OFF"
		log_queries_not_using_indexes = "OFF"  
		innodb_stats_on_metadata = "OFF"  
		table_open_cache_instances = 1  
		group_concat_max_len = 102
	}
}
`

const testAccKrdsParameterGroupUpdateConfig = `
provider "ksyun" {
	region = "cn-beijing-6"
}

resource "ksyun_krds_parameter_group" "foo" {
	name = "tf_test_on_hcl"
	engine = "mysql"
	engine_version = "5.7"
	description = "terraform created update"
	parameters = {
		connect_timeout = 60
		log_slow_admin_statements = "OFF"
		log_bin_trust_function_creators = "OFF"
		log_queries_not_using_indexes = "OFF"  
		innodb_stats_on_metadata = "OFF"  
		table_open_cache_instances = 1  
		group_concat_max_len = 102
		max_connect_errors = 2000
	}
}
`
const testAccKrdsParameterGroupApplyConfig = `

`
