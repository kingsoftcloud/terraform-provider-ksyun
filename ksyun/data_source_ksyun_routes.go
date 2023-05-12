/*
This data source provides a list of Route resources according to their Route ID, cidr and the VPC they belong to.

# Example Usage

```hcl

	data "ksyun_routes" "default" {
	  output_file="output_result"
	  ids=[]
	  vpc_ids=[]
	  instance_ids=[]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunRoutesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Route IDs, all the Route resources belong to this region will be retrieved if the ID is `\"\"`.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Route resources that satisfy the condition.",
			},

			"vpc_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of VPC id that the desired Route belongs to.",
			},

			"instance_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of the Route target id.",
			},

			"routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of routes. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of Route.",
						},

						"route_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of Route.",
						},

						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of VPC.",
						},

						"destination_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cidr block of the desired Route.",
						},

						"route_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the desired Route.",
						},

						"next_hop_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of next hop.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gateway_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the gateway.",
									},

									"gateway_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the gateway.",
									},
								},
							},
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time of the route.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunRoutesRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetRoutes(d, dataSourceKsyunRoutes())
}
