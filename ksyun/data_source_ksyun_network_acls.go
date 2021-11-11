package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunNetworkAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunNetworkAclsRead,

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
			"network_acls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"network_acl_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"network_acl_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"network_acl_entry_set": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_acl_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_acl_entry_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cidr_block": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rule_number": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"direction": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rule_action": {
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
func dataSourceKsyunNetworkAclsRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetNetworkAcls(d, dataSourceKsyunNetworkAcls())
}
