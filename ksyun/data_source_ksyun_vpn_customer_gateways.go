package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunVpnCustomerGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunVpnCustomerGatewaysRead,

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
			"customer_gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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

						"customer_gateway_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"customer_gateway_address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ha_customer_gateway_address": {
							Type:     schema.TypeString,
							Computed: true,
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
func dataSourceKsyunVpnCustomerGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetVpnCustomerGateways(d, dataSourceKsyunVpnCustomerGateways())
}
