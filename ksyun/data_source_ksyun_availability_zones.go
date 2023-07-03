/*
This data source provides a list of available zones in the current region.

# Example Usage

```hcl

	data "ksyun_availability_zones" "default" {
	  output_file=""
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunAvailabilityZonesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of AvailabilityZones that satisfy the condition.",
			},

			"availability_zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of AvailabilityZones. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the zone.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunAvailabilityZonesRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetAvailabilityZones(d, dataSourceKsyunAvailabilityZones())
}
