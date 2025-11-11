package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestClickhouse_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testClickhouseListConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_clickhouse.default"),
				),
			},
		},
	})
}

func TestKsyunClickhouse_details(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testClickhouseDetailsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_clickhouse.default"),
				),
			},
		},
	})
}

const testClickhouseListConfig = `

data "ksyun_clickhouse" "default"{
  output_file = "output_file"
}
`

const testClickhouseDetailsConfig = `
data "ksyun_clickhouse" "default"{
  output_file = "output_file"
  instance_id = "${clickhouse_instance_id}"
}
`
