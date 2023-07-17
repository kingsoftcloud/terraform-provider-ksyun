/*
This data source provides a list of VPN gateways.

# Example Usage

```hcl

	data "ksyun_vpn_gateways" "default" {
	  output_file="output_result"
	  ids=[]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunVpnGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunVpnGatewaysRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of VPN gateway IDs, all the resources belong to this region will be retrieved if the ID is `\"\"`.",
			},

			"vpc_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of VPC IDs.",
			},

			"project_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of project IDs.",
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by name.",
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
			"vpn_gateways": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the gateway.",
						},

						"vpn_gateway_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the gateway.",
						},

						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the gateway.",
						},

						"vpn_gateway_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the gateway.",
						},

						"band_width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Band width.",
						},

						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID.",
						},

						"gateway_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway IP address.",
						},

						"ha_gateway_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HA Gateway IP address.",
						},
						"vpn_gateway_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version of vpn gateway.",
						},

						"remote_cidr_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: "A list of remote cidrs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cidr block.",
									},
								},
							},
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
func dataSourceKsyunVpnGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetVpnGateways(d, dataSourceKsyunVpnGateways())
}
