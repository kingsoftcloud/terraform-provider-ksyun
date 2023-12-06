/*
This data source provides a list of VPN GatewayRoutes.

# Example Usage

```hcl

	data "ksyun_vpn_gateway_routes" "default" {
	  output_file="output_result"

	  # specify vpn_gateway_id to query vpn_gateway_routes
	  vpn_gateway_id = ""
	}

```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunVpnGatewayRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunVpnGatewayRoutesRead,

		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A list of VPN gateway IDs, all the resources belong to this region will be retrieved if the ID is `\"\"`.",
			},

			"next_hop_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of the next hop type.",
			},

			"cidr_blocks": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of cidr block.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of resources that satisfy the condition.",
			},
			"vpn_gateway_routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpn_gateway_route_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the gateway.",
						},

						"destination_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the gateway.",
						},

						"route_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the gateway.",
						},

						"next_hop_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the gateway.",
						},

						"next_hop_instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Band width.",
						},

						"vpn_gateway_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID.",
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
					},
				},
			},
		},
	}
}
func dataSourceKsyunVpnGatewayRoutesRead(d *schema.ResourceData, meta interface{}) error {
	vpnService := VpnSrv{meta.(*KsyunClient)}
	return vpnService.ReadAndSetVpnGatewayRoutes(d, dataSourceKsyunVpnGatewayRoutes())
}
