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
	public_key="ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC5MPEC4hjObq7u9uZY+IHwJE71wRMFPmYMo5Qc56Z9TG/5irXh0l8cu1+Qi9VmmJCeYMU8FrVoaBThjRIgKjuuAF9gKuYx8tWEsURO33F+s0u410PPgpOVyHM6yRO9QNM9iEVBQRk2T8cLdfZuKPQRQH+jyVMAFXpomcx7Q0Yt9rFkZIjC3wBw16MziaCPKVSyCn5SQ+mFNCqqn5+lVc5gXhWAJRwfnSFBNJhXNkEuPAFm1UeNW34Zi96SRA2msIuTBmxt+ZtczIw+MGh2/L8wrTUgXg6j9uD80ZwUKRWkmvapPHaqHZ0gPTvbL6lXtfwU4u2MQrvGSSThKugp1QZf+OC8/F9B4K0ehq6NCjaIscSMW33hO96un6kSvz5HWL0mwlJ+ZXvV6mpUY9X8HSJhOxaA4uX5RO/KqcQWhIku7LxvRBsK8O5EUwCYs6tMCQqpdn+edK5NM7PrO3j+IG6NlHD/JPlNh4pjyzK8oPbfTncYhtwBsAaWQLpW/oFxYv8= lvsongke@lvsongkedeMacBook-Air.local"
}
`

const testAccSSHKeyUpdateConfig = `
resource "ksyun_ssh_key" "foo" {
	key_name="sshKeyName-update"
}
`
