/*
This data source provides a list of KNAD resources  according to their KNAD ID.

# Example Usage

```hcl

	data "ksyun_knads" "default" {
	  output_file = "output_result"

	  project_id = []
	  ids    = []
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKnads() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKnadRead,
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
				Description: "A list of Knad IDs, all the Knads belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"knads": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of Knad. Each element contains the following attributes:",
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
							Description: "the bill type of the Knad. Valid Values: 1:(PrePaidByMonth),5:(DailyPaidByTransfer).",
						},
						"exprie_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the exprie time of the Knad.",
						},
						"idc_band": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the idcband of the Knad.",
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
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the service id of the Knad.Valid Values:'KNAD_30G','KNAD_100G','KNAD_300G','KNAD_1000G',''KNAD_2000G''.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKnadRead(d *schema.ResourceData, meta interface{}) error {
	knadService := KnadService{meta.(*KsyunClient)}
	return knadService.ReadAndSetKnads(d, dataSourceKsyunKnads())
}
