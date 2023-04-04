/*
This data source provides a list of VPN custom gateways.

# Example Usage

```hcl

	data "ksyun_vpn_customer_gateways" "default" {
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

func dataSourceKsyunVpnCustomerGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunVpnCustomerGatewaysRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of VPN customer gateway IDs, all the resources belong to this region will be retrieved if the ID is `\"\"`.",
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
			"customer_gateways": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the customer gateway.",
						},

						"customer_gateway_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the customer gateway.",
						},

						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the customer gateway.",
						},

						"customer_gateway_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the customer gateway.",
						},

						"customer_gateway_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the customer gateway.",
						},

						"ha_customer_gateway_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the HA customer gateway.",
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time.",
						},
					},
				},
			},
		},
	}
}
func dataSourceKsyunVpnCustomerGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetVpnCustomerGateways(d, dataSourceKsyunVpnCustomerGateways())
}
