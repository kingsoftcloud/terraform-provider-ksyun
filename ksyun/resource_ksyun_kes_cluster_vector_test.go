package ksyun

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// testAccKesClusterVectorRequireRealIDs fails the test fast (without contacting
// the API) when the VPC / subnet / security-group ids are missing or still the
// placeholder, or when the Xpack password is missing/placeholder. This prevents
// accidentally running the acceptance test with placeholder values, which would
// fail at the API at best.
func testAccKesClusterVectorRequireRealIDs(t *testing.T) {
	testAccKesClusterRequireRealIDs(t)
	if v := os.Getenv("KSYUN_KES_XPACK_PASSWORD"); v == "" || strings.HasPrefix(v, "your-") {
		t.Fatalf("KSYUN_KES_XPACK_PASSWORD must be set to a real password for the KES vector acceptance test")
	}
}

// TestAccKsyunKesClusterVector_basic verifies that a KES8 (vector) cluster
// resource can be created and read back. It requires real Ksyun credentials
// (KSYUN_ACCESS_KEY/KSYUN_SECRET_KEY), TF_ACC=1, real resource ids
// (KSYUN_KES_VPC_DOMAIN_ID/KSYUN_KES_VPC_SUBNET_ID/KSYUN_KES_SECURITY_GROUP_ID),
// and a real Xpack password (KSYUN_KES_XPACK_PASSWORD).
//
// NOTE: the KES SDK has no DeleteCluster API, so Delete is a no-op. A successful
// run leaves a real (billed) cluster in your account; delete it manually in the
// Ksyun console afterwards.
func TestAccKsyunKesClusterVector_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccKesClusterVectorRequireRealIDs(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKesClusterVectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKesClusterVectorConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKesClusterVectorExists("ksyun_kes_cluster_vector.example"),
					resource.TestCheckResourceAttrSet("ksyun_kes_cluster_vector.example", "cluster_id"),
					resource.TestCheckResourceAttrSet("ksyun_kes_cluster_vector.example", "cluster_name"),
				),
			},
		},
	})
}

// testAccKesClusterVectorConfig builds the HCL from environment variables. The
// cluster name and main version use vector-specific env vars
// (KSYUN_KES_VECTOR_*) so that running the KES7 and KES8 tests in the same
// shell does not let one override the other's version. The shared resource ids
// (VPC / subnet / security group) reuse the same env vars as the KES7 test.
//
// On the kmr/2021-09-02 backend the per-group placement (availability_zone /
// vpc_id / vpc_subnet_id) and the per-group index / charge_type live INSIDE each
// instance_groups entry; there is no top-level availability_zone or
// purchase_time, so neither is modelled here.
func testAccKesClusterVectorConfig() string {
	// placement is shared across groups; resolve once and reuse per group.
	az := kesTestEnvOr("KSYUN_KES_AVAILABILITY_ZONE", "cn-beijing-6a")
	vpcId := kesTestEnvOr("KSYUN_KES_VPC_DOMAIN_ID", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
	vpcSubnetId := kesTestEnvOr("KSYUN_KES_VPC_SUBNET_ID", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
	return fmt.Sprintf(`
resource "ksyun_kes_cluster_vector" "example" {
  cluster_name      = %q
  main_version      = %q
  distribution      = %q
  scenario          = "General"
  charge_type       = "Monthly"
  project_id        = %s
  vpc_domain_id     = %q
  vpc_subnet_id     = %q
  security_group_id = %q

  # required for any 8.x / 9.x cluster
  xpack_username = %q
  xpack_password = %q

  # KES8-specific switches (defaults shown).
  xpack_enabled             = true
  enable_best_delivery_mode = false
  enable_eip                = false
  instance_login_type       = "later"
  services                  = ["elasticsearch"]
  plugins                   = ["repository-ks3"]

  instance_groups {
    instance_group_type     = "MASTER"
    instance_type           = "ES.X6.4C8G"
    instance_type_code      = "ES.X6.4C8G"
    resource_type           = "KEC"
    instance_count          = 3
    instance_group_index    = 1
    availability_zone       = %q
    availability_zone_index = 1
    vpc_id                  = %q
    vpc_subnet_id           = %q
    charge_type             = "Monthly"
    volume_type             = "SSD3.0"
    volume_size             = 100
    volume_count            = 1
    system_disk_type        = "SSD3.0"
    system_disk_size        = 50
  }

  instance_groups {
    instance_group_type     = "DATA"
    instance_type           = "ES.X6.4C8G"
    instance_type_code      = "ES.X6.4C8G"
    resource_type           = "KEC"
    instance_count          = 3
    instance_group_index    = 2
    availability_zone       = %q
    availability_zone_index = 1
    vpc_id                  = %q
    vpc_subnet_id           = %q
    charge_type             = "Monthly"
    volume_type             = "SSD3.0"
    volume_size             = 500
    volume_count            = 1
    system_disk_type        = "SSD3.0"
    system_disk_size        = 50
  }
}
`,
		kesTestEnvOr("KSYUN_KES_VECTOR_CLUSTER_NAME", "tf-acc-kes-vector-cluster"),
		kesTestEnvOr("KSYUN_KES_VECTOR_MAIN_VERSION", "8.19.13"),
		kesTestEnvOr("KSYUN_KES_DISTRIBUTION", "ElasticSearch"),
		kesTestEnvOr("KSYUN_KES_PROJECT_ID", "0"),
		vpcId,
		vpcSubnetId,
		kesTestEnvOr("KSYUN_KES_SECURITY_GROUP_ID", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"),
		kesTestEnvOr("KSYUN_KES_XPACK_USERNAME", "elastic"),
		kesTestEnvOr("KSYUN_KES_XPACK_PASSWORD", "your-xpack-password"),
		// MASTER group placement
		az, vpcId, vpcSubnetId,
		// DATA group placement
		az, vpcId, vpcSubnetId,
	)
}

func testAccCheckKesClusterVectorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("KES vector cluster ID is not set")
		}
		return nil
	}
}

// testAccCheckKesClusterVectorDestroy verifies that the resource has been
// removed from Terraform state. Note: the KES SDK does not provide a
// DeleteCluster API, so Delete is a no-op that only removes the resource from
// state; the underlying cluster must be deleted manually in the Ksyun console.
func testAccCheckKesClusterVectorDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_kes_cluster_vector" {
			continue
		}
		if rs.Primary.ID != "" {
			return fmt.Errorf("KES vector cluster %s still present in state (delete is a no-op; remove it from the console)", rs.Primary.ID)
		}
	}
	return nil
}
