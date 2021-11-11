package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunNats() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunNatsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"vpc_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"project_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"nats": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"nat_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"nat_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"nat_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"nat_ip_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"band_width": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"nat_ip_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nat_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nat_ip_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"associate_nat_set": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunNatsRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetNats(d, dataSourceKsyunNats())
}
