package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSecurityGroupsRead,
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
			"vpc_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_entry_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_entry_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cidr_block": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"direction": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"icmp_type": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"icmp_code": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"port_range_from": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"port_range_to": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetSecurityGroups(d, dataSourceKsyunSecurityGroups())
}
