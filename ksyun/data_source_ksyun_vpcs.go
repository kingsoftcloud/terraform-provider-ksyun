/*
This data source provides a list of VPC resources according to their VPC ID, name.

# Example Usage

```hcl

	data "ksyun_vpcs" "default" {
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

func dataSourceKsyunVpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunVpcsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of VPC IDs, all the VPC resources belong to this region will be retrieved if the ID is `\"\"`.",
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by VPC name.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of VPC resources that satisfy the condition.",
			},
			"vpcs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of VPC.",
						},

						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of VPC.",
						},

						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of VPC.",
						},

						"vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of VPC.",
						},

						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CIDR blocks of VPC.",
						},
						"ipv6_cidr_block_association_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv6_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "the Ipv6 of this vpc bound.",
									},
								},
							},
							Description: "An Ipv6 association list of this vpc.",
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation for VPC.",
						},
					},
				},
			},
		},
	}
}
func dataSourceKsyunVpcsRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetVpcs(d, dataSourceKsyunVpcs())
}
