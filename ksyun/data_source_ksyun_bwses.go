/*
This data source provides a list of BWS resources (BandWidthShare) according to their BWS ID.

Example Usage

```hcl
data "ksyun_bwses" "default" {
  output_file="output_result"
  ids = ["c7b2ba05-9302-4933-8588-a66f920ff57d"]
}
```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunBandWidthShares() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunBandWidthSharesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of BWS IDs, all the BWSs belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"project_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more project IDs.",
			},

			"allocation_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more ids of the EIPs in the BWS.",
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by BWS name.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of BWS that satisfy the condition.",
			},
			"band_width_shares": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the BWS.",
						},

						"band_width_share_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the BWS.",
						},

						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of project.",
						},

						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the BWS.",
						},

						"band_width_share_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the BWS.",
						},

						"line_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the BWS line.",
						},

						"band_width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "bandwidth value.",
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time of the BWS.",
						},

						"allocation_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set:         schema.HashString,
							Description: "a list of EIP IDs which associated to BWS.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunBandWidthSharesRead(d *schema.ResourceData, meta interface{}) error {
	bwsService := BwsService{meta.(*KsyunClient)}
	return bwsService.ReadAndSetBandWidthShares(d, dataSourceKsyunBandWidthShares())
}
