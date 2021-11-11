package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunRoutesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
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

			"instance_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"route_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"destination_cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"route_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"next_hop_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gateway_id": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"gateway_name": {
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
					},
				},
			},
		},
	}
}

func dataSourceKsyunRoutesRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetRoutes(d, dataSourceKsyunRoutes())
}
