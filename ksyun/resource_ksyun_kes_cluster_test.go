package ksyun

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// kesTestEnvOr returns the env var value, or fallback when unset/empty. It lets
// the acceptance tests be driven by environment variables (real resource ids,
// cluster name, version, ...) so account-specific values are never hardcoded in
// the committed test sources.
func kesTestEnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// testAccKesClusterRequireRealIDs fails the test fast (without contacting the
// API) when the VPC / subnet / security-group ids are missing or still the
// placeholder. This prevents accidentally running the acceptance test with
// placeholder ids, which would fail at the API at best.
func testAccKesClusterRequireRealIDs(t *testing.T) {
	for _, k := range []string{"KSYUN_KES_VPC_DOMAIN_ID", "KSYUN_KES_VPC_SUBNET_ID", "KSYUN_KES_SECURITY_GROUP_ID"} {
		v := os.Getenv(k)
		if v == "" || strings.HasPrefix(v, "xxxx") {
			t.Fatalf("%s must be set to a real resource id for the KES acceptance test", k)
		}
	}
}

// TestAccKsyunKesClusterResource_basic verifies that a KES7 (standard) cluster
// resource can be created and read back. It requires real Ksyun credentials
// (KSYUN_ACCESS_KEY/KSYUN_SECRET_KEY), TF_ACC=1, and real resource ids
// (KSYUN_KES_VPC_DOMAIN_ID/KSYUN_KES_VPC_SUBNET_ID/KSYUN_KES_SECURITY_GROUP_ID).
// Named *Resource to avoid colliding with the homonymous data-source test
// TestAccKsyunKesCluster_basic.
//
// NOTE: the KES SDK has no DeleteCluster API, so Delete is a no-op. A successful
// run leaves a real (billed) cluster in your account; delete it manually in the
// Ksyun console afterwards.
func TestAccKsyunKesClusterResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccKesClusterRequireRealIDs(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKesClusterConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKesClusterExists("ksyun_kes_cluster.example"),
					resource.TestCheckResourceAttrSet("ksyun_kes_cluster.example", "cluster_id"),
					resource.TestCheckResourceAttrSet("ksyun_kes_cluster.example", "cluster_name"),
				),
			},
		},
	})
}

func testAccKesClusterConfig() string {
	return fmt.Sprintf(`
resource "ksyun_kes_cluster" "example" {
  cluster_name      = %q
  main_version      = %q
  scenario          = "general"
  availability_zone = %q
  charge_type       = "Monthly"
  purchase_time     = 1
  project_id        = %s
  vpc_domain_id     = %q
  vpc_subnet_id     = %q
  security_group_id = %q

  instance_groups {
    instance_group_type = "MASTER"
    instance_type       = "ES.I3.2C4G"
    resource_type       = "KEC"
    instance_count      = 3
    volume_type         = "SSD3.0"
    volume_size         = 100
    volume_count        = 1
    system_disk_type    = "SSD3.0"
    system_disk_size    = 50
  }

  instance_groups {
    instance_group_type = "DATA"
    instance_type       = "ES.D4.8C16G"
    resource_type       = "KEC"
    instance_count      = 3
    volume_type         = "SSD3.0"
    volume_size         = 500
    volume_count        = 1
    system_disk_type    = "SSD3.0"
    system_disk_size    = 50
  }
}
`,
		kesTestEnvOr("KSYUN_KES_CLUSTER_NAME", "tf-acc-kes-cluster"),
		kesTestEnvOr("KSYUN_KES_MAIN_VERSION", "7.10.0"),
		kesTestEnvOr("KSYUN_KES_AVAILABILITY_ZONE", "cn-beijing-6a"),
		kesTestEnvOr("KSYUN_KES_PROJECT_ID", "0"),
		kesTestEnvOr("KSYUN_KES_VPC_DOMAIN_ID", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"),
		kesTestEnvOr("KSYUN_KES_VPC_SUBNET_ID", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"),
		kesTestEnvOr("KSYUN_KES_SECURITY_GROUP_ID", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"),
	)
}

func testAccCheckKesClusterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("KES cluster ID is not set")
		}
		return nil
	}
}

// testAccCheckKesClusterDestroy verifies the resource was removed from Terraform
// state. The KES SDK has no DeleteCluster API, so Delete is a no-op that only
// drops the resource from state; the cloud cluster must be deleted manually in
// the Ksyun console.
func testAccCheckKesClusterDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_kes_cluster" {
			continue
		}
		if rs.Primary.ID != "" {
			return fmt.Errorf("KES cluster %s still in state after destroy (delete is a no-op; delete the cluster in the console)", rs.Primary.ID)
		}
	}
	return nil
}
