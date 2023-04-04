package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pkg/errors"
	"strings"
	"testing"
)

func TestAccKcsNode_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKcsNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKcsNodeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKcsNodeExists("ksyun_redis_instance_node.default"),
				),
			},
		},
	})
}

func testAccCheckKcsNodeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("kcs node is create failure")
		}
		return nil
	}
}

func testAccCheckKcsNodeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*KsyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_redis_instance_node" {
			instanceCheck := make(map[string]interface{})
			fmt.Println(rs.Primary.Attributes["cache_id"])
			instanceCheck["CacheId"] = rs.Primary.Attributes["cache_id"]
			ptr, err := client.kcsv2conn.DescribeCacheReadonlyNode(&instanceCheck)
			// Verify the error is what we want
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "cannot be found") {
					return nil
				}
				return err
			}

			if ptr != nil && len((*ptr)["Data"].([]interface{})) != 0 {
				return errors.New("delete node failure")
			}
		}
	}

	return nil
}

const testAccKcsNodeConfig = `
variable "available_zone" {
  default = "cn-beijing-6a"
}

variable "protocol" {
  default = "4.0"
}

resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_subnet" "foo" {
  subnet_name      	= "ksyun-subnet-tf"
  cidr_block 		= "10.7.0.0/21"
  subnet_type 		= "Reserve"
  dhcp_ip_from 		= "10.7.0.2"
  dhcp_ip_to 		= "10.7.0.253"
  vpc_id 			= "${ksyun_vpc.default.id}"
  gateway_ip 		= "10.7.0.1"
  dns1 				= "198.18.254.41"
  dns2 				= "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_redis_sec_group" "default" {
  available_zone 	= "${var.available_zone}"
  name 				= "testTerraform000"
  description 		= "testTerraform000"
}

resource "ksyun_redis_instance" "default" {
  available_zone        = "${var.available_zone}"
  name                  = "MyRedisInstance1101"
  mode                  = 2
  capacity              = 1
  net_type              = 2
  security_group_id     = "${ksyun_redis_sec_group.default.id}"
  vnet_id 				= "${ksyun_subnet.default.id}"
  vpc_id 				= "${ksyun_vpc.default.id}"
  bill_type             = 5
  duration              = ""
  duration_unit         = ""
  pass_word             = "Shiwo1101"
  iam_project_id        = "0"
  slave_num             = 0  
  protocol              = "${var.protocol}"
  reset_all_parameters  = false
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

resource "ksyun_redis_instance_node" "default" {
  cache_id          = "${ksyun_redis_instance.default.id}"
  available_zone    = "${var.available_zone}"
}
`
