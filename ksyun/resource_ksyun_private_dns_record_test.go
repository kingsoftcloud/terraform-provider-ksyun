package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunPrivateDnsRecord_basic(t *testing.T) {
	// var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_private_dns_record.foo",
		Providers:     testAccProviders,
		// CheckDestroy:  testAccCheckPrivateDnsDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsRecordConfig,
				// Check: resource.ComposeTestCheckFunc(
				// 	testAccCheckPrivateDnsExists("ksyun_private_dns_Record.foo", &val),
				// 	testAccCheckPrivateDnsAttributes(&val),
				// ),
			},
		},
	})
}
func TestAccKsyunPrivateDnsRecord_update(t *testing.T) {
	// var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_private_dns_record.foo",
		Providers:     testAccProviders,
		// CheckDestroy:  testAccCheckPrivateDnsDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsRecordConfig,
				// Check: resource.ComposeTestCheckFunc(
				// 	testAccCheckPrivateDnsExists("ksyun_private_dns_Record.foo", &val),
				// 	testAccCheckPrivateDnsAttributes(&val),
				// ),
			},
			{
				Config: testAccPrivateDnsRecordConfigUpdate,
				// Check: resource.ComposeTestCheckFunc(
				// 	testAccCheckPrivateDnsExists("ksyun_private_dns_Record.foo", &val),
				// 	testAccCheckPrivateDnsAttributes(&val),
				// ),
			},
		},
	})
}

const testAccPrivateDnsRecordConfig = `
resource "ksyun_private_dns_record" "foo" {
	record_name = "tf-pdns-record"
	record_ttl = 360
	zone_id = "a5ae6bf0-0ff4-472f-af9e-158afab99915"
	type = "CNAME"
	record_value = "tf-record.com"
	priority = 300
	weight = 300
	port = 300
}
`

const testAccPrivateDnsRecordConfigUpdate = `
resource "ksyun_private_dns_record" "foo" {
	record_name = "tf-pdns-record"
	record_ttl = 3600
	zone_id = "a5ae6bf0-0ff4-472f-af9e-158afab99915"
	type = "CNAME"
	record_value = "tf-record.com"
	priority = 300
	weight = 300
	port = 300
}
`
