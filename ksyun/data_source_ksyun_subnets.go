/*
This data source provides a list of Subnet resources according to their Subnet ID, name and the VPC they belong to.

# Example Usage

```hcl

	data "ksyun_subnets" "default" {
	  output_file="output_result"
	  ids=[]
	  vpc_id=[]
	  nat_id=[]
	  network_acl_id=[]
	  subnet_type=[]
	  availability_zone_name=[]
	}

```
*/
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
				Set:         schema.HashString,
				Description: "A list of Subnet IDs, all the Subnet resources belong to this region will be retrieved if the ID is `\"\"`.",
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				Description:  "A regex string to filter results by subnet name.",
			},

			"vpc_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The id of the VPC that the desired Subnet belongs to.",
			},

			"nat_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The id of the NAT that the desired Subnet associated to.",
			},

			"network_acl_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The id of the ACL that the desired Subnet associated to.",
			},

			"availability_zone_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The availability zone that the desired Subnet belongs to.",
			},

			"subnet_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "one or more subnet types.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Subnet resources that satisfy the condition.",
			},
			"subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of subnets. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the subnet.",
						},

						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the subnet.",
						},

						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the subnet.",
						},

						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CIDR block assigned to the subnet.",
						},

						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the VPC.",
						},

						"subnet_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the subnet.",
						},

						"dhcp_ip_from": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DHCP start IP.",
						},

						"dhcp_ip_to": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DHCP end IP.",
						},

						"gateway_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP of gateway.",
						},

						"dns1": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dns1 of the subnet.",
						},

						"dns2": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dns2 of the subnet.",
						},

						"network_acl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the ACL that the desired Subnet associated to.",
						},

						"nat_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the NAT that the desired Subnet associated to.",
						},

						"availability_zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone.",
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time of the subnet.",
						},
						"availble_i_p_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "number of available IPs.",
						},
						"ipv6_cidr_block_association_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv6_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "the Ipv6 of this vpc bound.",
									},
								},
							},
							Description: "An Ipv6 association list of this vpc.",
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
