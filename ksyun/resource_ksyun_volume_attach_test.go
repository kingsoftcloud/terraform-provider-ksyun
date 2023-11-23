package ksyun

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunVolumeAttach_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeAttachDestory,
		Steps: []resource.TestStep{
			{
				Config: attachConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeAttachExists("ksyun_volume_attach.foo"),
				),
			},
		},
	})
}

func testAccCheckVolumeAttachExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("attach id is empty")
		}
		return nil
	}
}

func testAccCheckVolumeAttachDestory(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_volume_attach" {
			continue
		}
		client := testAccProvider.Meta().(*KsyunClient)
		volume := make(map[string]interface{})
		volume["VolumeId.1"] = strings.Split(rs.Primary.ID, ":")[0]
		ptr, err := client.ebsconn.DescribeVolumes(&volume)
		if err != nil {
			return err
		}
		if ptr != nil {
			resp, ok := (*ptr)["Volumes"]
			if !ok {
				continue
			}
			l := resp.([]interface{})
			if len(l) == 0 {
				continue
			} else {
				if volume, ok := l[0].(map[string]interface{}); ok {
					if volume["VolumeStatus"] == "in-use" {
						return fmt.Errorf("VolumeAttach still exist")
					}
				}
			}
		}
	}
	return nil
}

const testAccVolumeAttachConfig = `
data "ksyun_images" "centos-7_5" {
  output_file=""
  platform= "centos-7.5"
  is_public=true
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "default" {
  subnet_name      = "ksyun-subnet-tf"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Normal"
  dhcp_ip_from = "10.7.0.2"
  dhcp_ip_to = "10.7.0.253"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "cn-beijing-6a"
}
resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="ksyun-security-group"
}
resource "ksyun_instance" "default" {
  image_id="${data.ksyun_images.centos-7_5.images.0.image_id}"
  instance_type="S4.1A"
  system_disk{
    disk_type="SSD3.0"
    disk_size=30
  }
  data_disk_gb=0
  #max_count=1
  #min_count=1
  subnet_id="${ksyun_subnet.default.id}"
  instance_password="Xuan663222"
  keep_image_login=false
  charge_type="Daily"
  purchase_time=1
  security_group_id=["${ksyun_security_group.default.id}"]
  private_ip_address=""
  instance_name="ksyun-kec-tf"
  instance_name_suffix=""
  sriov_net_support=false
  project_id=0
  data_guard_id=""
  key_id=[]
  force_delete=true
}

resource "ksyun_volume" "default" {
  volume_name="ksyun_volume_tf_test"
  volume_type="SSD3.0"
  size=10
  charge_type="Daily"
  availability_zone="cn-beijing-6a"
  volume_desc="ksyun_volume_tf_test"
}

resource "ksyun_volume_attach" "foo" {
volume_id="${ksyun_volume.default.id}"
instance_id="${ksyun_instance.default.id}"
}
`

const attachConfig = `
resource "ksyun_volume" "default" {
  volume_name="ksyun_volume_tf_test"
  volume_type="SSD3.0"
  size=10
  charge_type="Daily"
  availability_zone="cn-beijing-6a"
  volume_desc="ksyun_volume_tf_test"
}

resource "ksyun_volume_attach" "foo" {
volume_id=ksyun_volume.default.id
instance_id="3f0c35a4-2294-4152-a864-676c4913c45a"
}
`
