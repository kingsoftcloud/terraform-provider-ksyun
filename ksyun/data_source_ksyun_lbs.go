/*
This data source provides a list of Load Balancer resources according to their Load Balancer ID, VPC ID and Subnet ID.

# Example Usage

```hcl

	data "ksyun_lbs" "default" {
	  output_file="output_result"
	  name_regex=""
	  ids=[]
	  state=""
	  vpc_id=[]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunLbs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunLbsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Load Balancer IDs, all the LBs belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter resulting lbs by name.",
			},

			"vpc_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The ID of the VPC linked to the Load Balancers.",
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"associate",
					"disassociate",
				}, false),
				Description: "state of the LB.",
			},

			"project_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ID of the project.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Load Balancers that satisfy the condition.",
			},

			"lbs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "public ip address.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the VPC linked to the Load Balancers.",
						},
						"load_balancer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the Load Balancer.",
						},
						"load_balancer_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Load Balancer.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the Load Balancer.",
						},
						"load_balancer_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "start or stop.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "associate or disassociate.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the subnet.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the project.",
						},
						"listeners_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the listeners.",
						},
						"ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP version.",
						},
						"is_waf": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether it is a waf LB or not.",
						},
						"lb_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the LB.",
						},
						"lb_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "status of the LB.",
						},
						"access_logs_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether accessLogs is enabled or not.",
						},
						"access_logs_s3_bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bucket for storing access logs.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunLbsRead(d *schema.ResourceData, meta interface{}) error {
	slbService := SlbService{meta.(*KsyunClient)}
	return slbService.ReadAndSetLoadBalancers(d, dataSourceKsyunLbs())
}
