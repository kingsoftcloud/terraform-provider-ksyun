package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"log"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"ksyun": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("KSYUN_ACCESS_KEY"); v == "" {
		t.Fatal("KSYUN_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("KSYUN_SECRET_KEY"); v == "" {
		t.Fatal("KSYUN_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("KSYUN_REGION"); v == "" {
		log.Println("[INFO] Test: Using cn-beijing-6 as test region")
		os.Setenv("KSYUN_REGION", "cn-beijing-6")
	}

	// 这个默认影响大多数产品线的test，先去掉
	//if v := os.Getenv("KSYUN_DOMAIN"); v == "" {
	//	os.Setenv("KSYUN_DOMAIN", DEFAULT_ENDPOINT)
	//}
}

func testAccCheckIDExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf(" Can't find resource or data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID is not be set")
		}
		return nil
	}
}
