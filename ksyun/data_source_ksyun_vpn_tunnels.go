package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunVpnTunnels() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunVpnTunnelsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"vpn_gateway_ids": {
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
			"vpn_tunnels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpn_tunnel_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpn_gre_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"customer_gre_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ha_vpn_gre_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ha_customer_gre_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpn_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"customer_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpn_tunnel_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"pre_shared_key": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ike_authen_algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ike_dh_group": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"ike_encry_algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ipsec_encry_algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ipsec_authen_algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ipsec_life_time_traffic": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"ipsec_life_time_second": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"extra_cidr_set": {
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
func dataSourceKsyunVpnTunnelsRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetVpnTunnels(d, dataSourceKsyunVpnTunnels())
}
