/*
This data source provides a list of Subnet resources according to their Subnet ID, name and the VPC they belong to.

# Example Usage

```hcl

	data "ksyun_cens" "default" {
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

func dataSourceKsyunCens() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunCensRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Cen IDs, all the Cens belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by cen name.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of cens that satisfy the condition.",
			},
			"cens": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of cens. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the cen.",
						},

						"cen_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the cen.",
						},

						"cen_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the cen.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of cen.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time of the cen.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunCensRead(d *schema.ResourceData, meta interface{}) error {
	cenService := CenService{meta.(*KsyunClient)}
	return cenService.ReadAndSetCens(d, dataSourceKsyunCens())
}
