package ksyun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKsyunSSHKey_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_ssh_key.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSSHKeyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSSHKeyExists("ksyun_ssh_key.foo", &val),
					testAccCheckSSHKeyAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunSSHKey_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_ssh_key.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSSHKeyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSSHKeyExists("ksyun_ssh_key.foo", &val),
					testAccCheckSSHKeyAttributes(&val),
				),
			},
			{
				Config: testAccSSHKeyUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSSHKeyExists("ksyun_ssh_key.foo", &val),
					testAccCheckSSHKeyAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckSSHKeyExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("SSHKey id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		sSHKey := make(map[string]interface{})
		sSHKey["KeyId.1"] = rs.Primary.ID
		ptr, err := client.sksconn.DescribeKeys(&sSHKey)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["KeySet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}
		*val = *ptr
		return nil
	}
}
func testAccCheckSSHKeyAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["KeySet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("SSHKey id is empty")
			}
		}
		return nil
	}
}
func testAccCheckSSHKeyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_ssh_key" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		sSHKey := make(map[string]interface{})
		sSHKey["KeyId.1"] = rs.Primary.ID
		ptr, err := client.sksconn.DescribeKeys(&sSHKey)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["KeySet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("SSHKey still exist")
			}
		}
	}

	return nil
}

const testAccSSHKeyConfig = `
resource "ksyun_ssh_key" "foo" {
	key_name="sshKeyName"
}
`

const testAccSSHKeyUpdateConfig = `
resource "ksyun_ssh_key" "foo" {
	key_name="sshKeyName-update"
}
`
