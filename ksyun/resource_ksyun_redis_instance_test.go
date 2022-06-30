package ksyun

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pkg/errors"
	"testing"
)

func TestAccKcs_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_redis_instance.default",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKcsDestroy,
		Steps: []resource.TestStep{
			// 集群创建
			{
				Config: testAccKcsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKcsInstanceExists("ksyun_redis_instance.default"),
				),
			},
			// 集群更配
			{
				Config: testUpdateAccKcsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKcsInstanceExists("ksyun_redis_instance.default"),
				),
			},
		},
	})

	//resource.Test(t, resource.TestCase{
	//	PreCheck: func() {
	//		testAccPreCheck(t)
	//	},
	//	IDRefreshName: "ksyun_redis_instance.single",
	//	Providers:     testAccProviders,
	//	CheckDestroy:  testAccCheckKcsDestroy,
	//	Steps: []resource.TestStep{
	//		// 主从创建
	//		{
	//			Config: testAccKcsSingleConfig,
	//			Check: resource.ComposeTestCheckFunc(
	//				testAccCheckKcsInstanceExists("ksyun_redis_instance.single"),
	//			),
	//		},
	//		// 主从更配
	//		{
	//			Config: testAccKcsSingleUpdateConfig,
	//			Check: resource.ComposeTestCheckFunc(
	//				testAccCheckKcsInstanceExists("ksyun_redis_instance.single"),
	//			),
	//		},
	//	},
	//})
}

func testAccCheckKcsInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("kcs instance is create failure")
		}
		return nil
	}
}

func testAccCheckKcsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*KsyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_redis_instance" {
			instanceCheck := make(map[string]interface{})
			instanceCheck["CacheId"] = rs.Primary.ID
			ptr, err := client.kcsv1conn.DescribeCacheCluster(&instanceCheck)
			// Verify the error is what we want
			if err != nil {
				if ksyunError, ok := err.(awserr.RequestFailure); ok && ksyunError.StatusCode() == 404 {
					return nil
				}
				return err
			}
			if ptr != nil {
				if (*ptr)["Data"] != nil {
					if v, ok := (*ptr)["Data"].(map[string]interface{}); ok {
						orderUse := int(v["orderUse"].(float64))
						if orderUse == 7 {
							return nil
						}
					}
					return errors.New("delete instance failure")
				}
			}
		}
	}

	return nil
}

const testAccKcsConfig = `
variable "protocol" {
  default = "4.0"
}

resource "ksyun_redis_instance" "default" {
  available_zone        = "eu-east-1b"
  name                  = "HoubinTerraformTest"
  mode                  = 3
  capacity              = 6
  net_type              = 2
  security_group_id     = "20b7a66e-c2c5-4ccc-a260-b92ce8992adc"
  vnet_id 				= "b6a38f5a-804d-4a12-ab60-fa3e9681b3ae"
  vpc_id 				= "8c8a161f-bc75-46ea-8215-8db972a6cb51"
  bill_type             = 5
  duration              = ""
  pass_word             = "Shiwo1101"
  iam_project_id        = "0"
  slave_num             = 0
  protocol              = "${var.protocol}"
  reset_all_parameters  = false
  timing_switch         = "On"
  timezone              = "06:00-07:00"
  shard_size            = 2
  shard_num             = 3
  parameters = {
    "appendonly"                  = "no",
    "appendfsync"                 = "everysec",
    "maxmemory-policy"            = "volatile-lru",
    "hash-max-ziplist-entries"    = "512",
    "zset-max-ziplist-entries"    = "128",
    "list-max-ziplist-size"       = "-2",
    "hash-max-ziplist-value"      = "64",
    "notify-keyspace-events"      = "",
    "zset-max-ziplist-value"      = "64",
    "maxmemory-samples"           = "5",
    "set-max-intset-entries"      = "512",
    "timeout"                     = "600",
  }
}

resource "ksyun_tag" "test_tag" {
    key = "exist_tag"
    value = "exist_tag_value1"
    resource_type = "kcs-instance"
    resource_id = "${ksyun_redis_instance.default.id}"
}
`
const testUpdateAccKcsConfig = `
variable "protocol" {
  default = "4.0"
}

resource "ksyun_redis_instance" "default" {
  available_zone        = "eu-east-1b"
  name                  = "HoubinTerraformTest"
  mode                  = 3
  capacity              = 10
  net_type              = 2
  security_group_id     = "20b7a66e-c2c5-4ccc-a260-b92ce8992adc"
  vnet_id 				= "b6a38f5a-804d-4a12-ab60-fa3e9681b3ae"
  vpc_id 				= "8c8a161f-bc75-46ea-8215-8db972a6cb51"
  bill_type             = 5
  duration              = ""
  pass_word             = "wwsNewPwd123"
  iam_project_id        = "0"
  slave_num             = 0
  protocol              = "${var.protocol}"
  reset_all_parameters  = false
  timing_switch         = "On"
  timezone              = "07:00-08:00"
  shard_size            = 2
  shard_num             = 5
  parameters = {
    "appendonly"                  = "no",
    "appendfsync"                 = "everysec",
    "maxmemory-policy"            = "volatile-lru",
    "hash-max-ziplist-entries"    = "512",
    "zset-max-ziplist-entries"    = "128",
    "list-max-ziplist-size"       = "-2",
    "hash-max-ziplist-value"      = "64",
    "notify-keyspace-events"      = "",
    "zset-max-ziplist-value"      = "64",
    "maxmemory-samples"           = "5",
    "set-max-intset-entries"      = "512",
    "timeout"                     = "600",
  }
}

resource "ksyun_tag" "test_tag" {
    key = "exist_tag"
    value = "exist_tag_value1"
    resource_type = "kcs-instance"
    resource_id = "${ksyun_redis_instance.default.id}"
}
`

const testAccKcsSingleConfig = `
variable "protocol" {
  default = "4.0"
}

resource "ksyun_redis_instance" "single" {
  name                  = "TerraformRedisSingleTest"
  mode                  = 2
  capacity              = 1
  net_type              = 2
  security_group_id     = "20b7a66e-c2c5-4ccc-a260-b92ce8992adc"
  vnet_id 				= "74274e94-5e3e-4146-b962-92c0a59bb4e7"
  vpc_id 				= "c7d1a637-657d-4250-9c3d-d2426cc24de2"
  bill_type             = 5
  duration              = ""
  pass_word             = "Shiwo1101"
  iam_project_id        = "0"
  slave_num             = 1
  protocol              = "${var.protocol}"
  reset_all_parameters  = false
  timing_switch         = "On"
  timezone              = "06:00-07:00"
  available_zone        = "cn-beijing-6b"
  prepare_az_name       = "cn-beijing-6a"
  rr_az_name            = "cn-beijing-6b"
  parameters = {
    "appendonly"                  = "no",
    "appendfsync"                 = "everysec",
    "maxmemory-policy"            = "volatile-lru",
    "hash-max-ziplist-entries"    = "512",
    "zset-max-ziplist-entries"    = "128",
    "list-max-ziplist-size"       = "-2",
    "hash-max-ziplist-value"      = "64",
    "notify-keyspace-events"      = "",
    "zset-max-ziplist-value"      = "64",
    "maxmemory-samples"           = "5",
    "set-max-intset-entries"      = "512",
    "timeout"                     = "600",
  }
}

resource "ksyun_tag" "test_tag" {
    key = "exist_tag"
    value = "exist_tag_value2"
    resource_type = "kcs-instance"
    resource_id = "${ksyun_redis_instance.single.id}"
}
`

const testAccKcsSingleUpdateConfig = `
variable "protocol" {
  default = "4.0"
}

resource "ksyun_redis_instance" "single" {
  name                  = "TerraformRedisSingleTest123"
  mode                  = 2
  capacity              = 2
  net_type              = 2
  security_group_id     = "20b7a66e-c2c5-4ccc-a260-b92ce8992adc"
  vnet_id 				= "74274e94-5e3e-4146-b962-92c0a59bb4e7"
  vpc_id 				= "c7d1a637-657d-4250-9c3d-d2426cc24de2"
  bill_type             = 5
  duration              = ""
  pass_word             = "Shiwo1101xxxx"
  iam_project_id        = "0"
  slave_num             = 1
  protocol              = "${var.protocol}"
  reset_all_parameters  = false
  timing_switch         = "On"
  timezone              = "07:00-08:00"
  available_zone        = "cn-beijing-6b"
  prepare_az_name       = "cn-beijing-6a"
  rr_az_name            = "cn-beijing-6b"
  parameters = {
    "appendonly"                  = "no",
    "appendfsync"                 = "everysec",
    "maxmemory-policy"            = "volatile-lru",
    "hash-max-ziplist-entries"    = "512",
    "zset-max-ziplist-entries"    = "128",
    "list-max-ziplist-size"       = "-2",
    "hash-max-ziplist-value"      = "64",
    "notify-keyspace-events"      = "",
    "zset-max-ziplist-value"      = "64",
    "maxmemory-samples"           = "5",
    "set-max-intset-entries"      = "512",
    "timeout"                     = "600",
  }
}

resource "ksyun_tag" "test_tag" {
    key = "exist_tag"
    value = "exist_tag_value3"
    resource_type = "kcs-instance"
    resource_id = "${ksyun_redis_instance.single.id}"
}
`
