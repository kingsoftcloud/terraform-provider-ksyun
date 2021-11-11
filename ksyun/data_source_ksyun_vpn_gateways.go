package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunVpnGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunVpnGatewaysRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
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
			"vpn_gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpn_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpn_gateway_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"band_width": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"gateway_address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ha_gateway_address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"remote_cidr_set": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_block": {
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
func dataSourceKsyunVpnGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetVpnGateways(d, dataSourceKsyunVpnGateways())
}
