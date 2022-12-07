/*
This data source provides a list of available zones in the current region.

# Example Usage

```hcl

	data "ksyun_availability_zones" "default" {
	  output_file=""
	  ids=[]
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
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"availability_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the zone",
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
