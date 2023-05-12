/*
This data source provides a list of Security Group resources according to their Security Group ID, name and resource id.

# Example Usage

```hcl

	data "ksyun_security_groups" "default" {
	  output_file="output_result"
	  ids=[]
	  vpc_id=[]
	}

```
*/
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
				Set:         schema.HashString,
				Description: "A list of Security Group IDs, all the Security Group resources belong to this region will be retrieved if the ID is `\"\"`.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Security Group resources that satisfy the condition.",
			},
			"vpc_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of VPC IDs.",
			},
			"security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of security groups. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation for the security group.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the VPC.",
						},
						"security_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the security group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the security group.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the security group.",
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the security group.",
						},
						"security_group_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the security group.",
						},
						"security_group_entry_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of the security group entries.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of the security group entry.",
									},
									"security_group_entry_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the security group entry.",
									},
									"cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The cidr block of source.",
									},
									"direction": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The direction of the entry.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "protocol of the entry.",
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
										Description: "The start of port numbers.",
									},
									"port_range_to": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The end of port numbers.",
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
