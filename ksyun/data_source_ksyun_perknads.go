/*
This data source provides a list of PerPay KNAD resources according to their KNAD ID.

# Example Usage

```hcl

	data "ksyun_perknads" "default" {
	  output_file = "output_result"

	  project_id = []
	  ids        = []
	}

```
*/
package ksyun

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataSourceKsyunPerKnads() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunPerKnadsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more project IDs.",
			},
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Knad IDs, all the PerKnads belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"perknads": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of PerPay Knad. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"knad_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the resource.",
						},
						"band": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the band of the Knad.",
						},
						"bill_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the bill type of the Knad.",
						},
						"exprie_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the exprie time of the Knad.",
						},
						"ip_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the max ip count that can bind to the Knad.",
						},
						"knad_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the name of the Knad.",
						},
						"max_band": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the max band of the Knad.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the project.",
						},
						"used_ip_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The binding ip count of the Knad.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunPerKnadsRead(d *schema.ResourceData, meta interface{}) error {
	svc := PerKnadService{meta.(*KsyunClient)}
	return svc.ReadAndSetPerKnads(d, dataSourceKsyunPerKnads())
}
