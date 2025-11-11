/*
This data source provides a list of Direct Connect resources.

Example Usage

```hcl
data "ksyun_direct_connects" "test" {
  ids        = []
  name_regex = ".*test.*"
}


```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunDirectConnects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunDirectConnectsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Direct Connect IDs.",
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by Direct Connect name.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Direct Connect that satisfy the condition.",
			},
			"direct_connects": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the Direct Connect.",
						},

						"direct_connect_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the Direct Connect.",
						},

						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type.",
						},

						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the Direct Connect.",
						},

						"direct_connect_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the Direct Connect.",
						},

						"pop_location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pop Location.",
						},

						"customer_location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Customer location of the Direct Connect.",
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time of the Direct Connect.",
						},

						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State.",
						},
						"band_width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Band Width.",
						},

						"vlan": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Vlan.",
						},
						"distance": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Distance.",
						},
						"vpc_noc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vpc Noc ID.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunDirectConnectsRead(d *schema.ResourceData, meta interface{}) error {
	srv := VpcService{meta.(*KsyunClient)}
	return srv.ReadAndSetDirectConnects(d, dataSourceKsyunDirectConnects())
}
