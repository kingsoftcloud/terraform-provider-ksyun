/*
This data source provides a list of Private Dns Zone.

# Example Usage

```hcl

data "ksyun_private_dns_zones" "foo" {
  output_file = "pdns_output_result"
  zone_ids = []
}
```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunPrivateDnsZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunPrivateDnsZonesRead,
		Schema: map[string]*schema.Schema{
			"zone_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.Any(
						validation.StringIsNotWhiteSpace,
					),
				},
				Set:         schema.HashString,
				Description: "A list of the filter values that is zone id.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Private Dns Zone that satisfy the condition.",
			},
			"zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of Private Dns Zone. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the Private Dns Zone.",
						},
						"zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of Private Dns Zone.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project Id.",
						},
						"zone_ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Zone TTL.",
						},
						"bind_vpc_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The zone Bound VPCs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region Name.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC id.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status of Zone.",
									},
									"vpc_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The VPC name.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunPrivateDnsZonesRead(d *schema.ResourceData, meta interface{}) error {
	s := DnsService{meta.(*KsyunClient)}
	return s.ReadAndSetPrivateDnsZones(d, dataSourceKsyunPrivateDnsZones())
}
