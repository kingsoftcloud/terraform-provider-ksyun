// ksyun-terraform
terraform {
  required_providers {
    ksyun = {
        source = "kingsoftcloud/ksyun"
    }
  }
}

provider "ksyun" {
  region = "cn-beijing-6"
}

resource "ksyun_krds_parameter_group" "dpg2" {
    name = "tf_dpg_on_hcl"
    description = "tf configuration test"
    engine = "mysql"
    engine_version = "5.6"
    parameters = {
        connect_timeout = 20
        innodb_stats_on_metadata = "OFF"  
		table_open_cache_instances = 1  
		group_concat_max_len = 102
		max_connect_errors = 2000
        max_prepared_stmt_count = 65535
        max_user_connections = 65535
    }
}

resource "ksyun_krds_parameter_group" "dpg1" {
    name = "tf_test_post"
    # description = "tf configuration test"
    engine = "mysql"
    engine_version = "5.6"
    parameters = {
        connect_timeout = "30"
    }
}

data "ksyun_krds_parameter_group" "foo" {
	output_file = "output_result"
	db_parameter_group_id = ksyun_krds_parameter_group.dpg1.id
}

output "dpg_out" {
    value = data.ksyun_krds_parameter_group.foo
}
