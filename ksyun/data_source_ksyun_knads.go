package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKnads() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKnadRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"knad_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"duration": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("bill_type"); ok && v == 1 {
						return false
					}
					return true
				},
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"band": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_band": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"ip_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"bill_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"idc_band": {
				Type:     schema.TypeInt,
				Required: true,
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
						"charge_mode": {
							Type:     schema.TypeString,
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
