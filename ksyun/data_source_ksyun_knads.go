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
				Set: schema.HashString,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"knads": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"knad_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"band": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bill_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"exprie_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"idc_band": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"ip_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"knad_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_band": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
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
