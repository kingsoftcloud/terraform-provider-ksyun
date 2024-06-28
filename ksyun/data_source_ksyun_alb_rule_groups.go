/*
This data source provides a list of ALB rule group resources according to their ID.

# Example Usage

```hcl

	data "ksyun_alb_rule_groups" "default" {
		output_file="output_result"
		ids=[]
		alb_listener_id=[]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunAlbRuleGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunAlbRuleGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of ALB Rule Group IDs, all the ALB Rule Group belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"alb_listener_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "one or more alb listener id.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of ALB Rule Groups that satisfy the condition.",
			},
			"alb_rule_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of ALB Rule Groups. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the ALB Rule Group.",
						},
						"alb_listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the ALB listener.",
						},
						"alb_rule_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the rule group.",
						},
						"alb_rule_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the ALB rule group.",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Forwarding mode of listener.",
						},
						"backend_server_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the backend server group.",
						},
						"listener_sync": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to synchronize the health check, session persistence, and load balancing algorithm of the listener.",
						},
						"session_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of session.",
						},
						"session_persistence_period": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Session hold timeout.",
						},
						"cookie_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of cookie.",
						},
						"cookie_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of cookie.",
						},
						"health_check_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status maintained by health examination.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Interval of health examination.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Health check timeout.",
						},
						"healthy_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Health threshold.",
						},
						"unhealthy_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unhealthy threshold.",
						},
						"url_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to HTTP type listener health check.",
						},
						"host_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The service host name of the health check.",
						},
						"http_method": {
							Type:        schema.TypeString,
							Description: "The http requests' method.",
							Computed:    true,
						},

						"health_port": {
							Type:        schema.TypeInt,
							Description: "The port of connecting for health check.",
							Computed:    true,
						},

						"health_protocol": {
							Type:        schema.TypeString,
							Description: "The protocol of connecting for health check.",
							Computed:    true,
						},
						"alb_rule_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Rule set.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alb_rule_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule type.",
									},
									"alb_rule_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule value.",
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

func dataSourceKsyunAlbRuleGroupsRead(d *schema.ResourceData, meta interface{}) error {
	s := AlbRuleGroup{meta.(*KsyunClient)}
	return s.ReadAndSetRuleGroups(d, dataSourceKsyunAlbRuleGroups())
}
