/*
This data source provides a list of ALB resources according to their ALB ID.

# Example Usage

```hcl

	data "ksyun_albs" "default" {
		output_file="output_result"
		ids=[]
		vpc_id=[]
		state=[]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunAlbs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunAlbsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of ALB IDs, all the ALBs belong to this region will be retrieved if the ID is `\"\"`.",
			},
			// openapi暂未支持
			//"project_id": {
			//	Type:     schema.TypeSet,
			//	Optional: true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeString,
			//	},
			//	Set:         schema.HashString,
			//	Description: "One or more project IDs.",
			//},
			"vpc_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more VPC IDs.",
			},
			"state": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more state.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of ALBs that satisfy the condition.",
			},
			"albs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of ALB. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the ALB.",
						},
						"alb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the ALB.",
						},
						"alb_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the ALB.",
						},
						"alb_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version of the ALB.",
						},
						"alb_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the ALB.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the VPC.",
						},
						"ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP version, 'ipv4' or 'ipv6'.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the project.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The charge type.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time.",
						},
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public IP address.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the ALB.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the ALB.",
						},
						"enabled_log": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether log is enabled or not.",
						},
						"klog_info": {
							Type:        schema.TypeList,
							Computed:    true,
							MaxItems:    1,
							Description: "klog info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "account id.",
									},
									"log_pool_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "log pool name.",
									},
									"project_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "log project name.",
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

func dataSourceKsyunAlbsRead(d *schema.ResourceData, meta interface{}) error {
	s := AlbService{meta.(*KsyunClient)}
	return s.ReadAndSetAlbs(d, dataSourceKsyunAlbs())
}
