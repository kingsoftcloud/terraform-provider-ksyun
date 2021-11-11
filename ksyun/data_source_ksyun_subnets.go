package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSubnetsRead,
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
				ValidateFunc: validation.ValidateRegexp,
			},

			"vpc_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"nat_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"network_acl_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"availability_zone_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"subnet_types": {
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
			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"subnet_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"dhcp_ip_from": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"dhcp_ip_to": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"dns1": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"dns2": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"network_acl_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"nat_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"availability_zone_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availble_i_p_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetSubnets(d, dataSourceKsyunSubnets())
}
