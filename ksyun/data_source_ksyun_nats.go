/*
This data source provides a list of Nat resources according to their Nat ID and the VPC they belong to.

# Example Usage

```hcl

	data "ksyun_nats" "default" {
	  output_file="output_result"
	  ids=[]
	  vpc_ids=[]
	  project_ids=[]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunNats() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunNatsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Nat IDs, all the Nat resources belong to this region will be retrieved if the ID is `\"\"`.",
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by NAT name.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of NAT resources that satisfy the condition.",
			},

			"vpc_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of VPC id that the desired Nat belongs to.",
			},

			"project_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Project id that the desired Nat belongs to.",
			},

			"nats": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of NATs. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of NAT.",
						},

						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC ID of the desired Nat belongs to.",
						},

						"nat_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the NAT.",
						},

						"nat_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The mode of the NAT.",
						},

						"nat_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the NAT.",
						},

						"nat_ip_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The nat ip count of the desired Nat.",
						},

						"band_width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The nat ip band width of the desired Nat.",
						},

						"nat_ip_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The nat ip list of the desired Nat.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nat_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "NAT IP address.",
									},
									"nat_ip_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the NAT IP.",
									},
								},
							},
						},

						"associate_nat_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: "The subnet associate list of the desired Nat.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the subnet.",
									},
								},
							},
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation of Nat.",
						},

						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the project.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunNatsRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetNats(d, dataSourceKsyunNats())
}
