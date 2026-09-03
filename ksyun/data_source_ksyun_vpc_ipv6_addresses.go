/*
This data source provides a list of line resources supported.

# Example Usage

```hcl

	data "ksyun_vpc_ipv6_addresses" "default" {
	  ids=["54098712-6a7b-4561-a84a-d36e9bb0b206"]
	  network_interface_id=["54098712-6a7b-4561-a84a-d36e9bb0b206"]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunVpcIpv6Addresses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunVpcIpv6AddressesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more ipv6 public ip address IDs.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by Bare Metal name.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of lines that satisfy the condition.",
			},

			"network_interface_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "Network interface id.",
			},

			"ipv6_public_ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All the  Ipv6PublicIpAddress according the argument.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"band_width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Bandwidth.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Charge type.",
						},
						"service_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service end time.",
						},
						"ipv6_public_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Ipv6 public ip address.",
						},
						"ipv6_public_ip_address_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Ipv6 public ip address id.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunVpcIpv6AddressesRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetVpcIpv6Addresses(d, dataSourceKsyunVpcIpv6Addresses())
}
