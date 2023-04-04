/*
This data source provides a list of Network ACL resources

# Example Usage

```hcl

	data "ksyun_network_acls" "default" {
	  output_file="output_result"

	//  vpc_ids = ["769c780b-acbd-41ca-9a06-4960e2423c7e"]
	}

```
*/
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
				Set:         schema.HashString,
				Description: "A list of network ACL IDs.",
			},

			"vpc_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of VPC IDs.",
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by ACL name.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of ACLs that satisfy the condition.",
			},
			"network_acls": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of ACLs. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the ACL.",
						},

						"network_acl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the ACL.",
						},

						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the VPC.",
						},

						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the ACL.",
						},

						"network_acl_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the ACL.",
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time of the ACL.",
						},

						"network_acl_entry_set": {
							Type:     schema.TypeList,
							Computed: true,
							//Optional:    true,
							Description: "A list of the ACL entries.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the ACL entry.",
									},
									"network_acl_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the ACL.",
									},
									"network_acl_entry_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the ACL entry.",
									},
									"cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The information of Acl's cidr block.",
									},
									"rule_number": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "rule priority.",
									},
									"direction": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "rule direction.",
									},
									"rule_action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "rule action, allow or deny.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "rule protocol.",
									},
									"icmp_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ICMP type.",
									},
									"icmp_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ICMP code.",
									},
									"port_range_from": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "beginning of the port range.",
									},
									"port_range_to": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ending of the port range.",
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
