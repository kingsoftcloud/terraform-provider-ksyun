/*
Provides a list of ksyun lb rules resources in the current region.

# Example Usage

```

	data "ksyun_lb_rules" "default" {
		output_file="output_result"
		ids=[]
		host_header_id=[]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunSlbRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSlbRulesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of rule IDs.",
			},
			"host_header_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The id of host header.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running terraform plan).",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of LB Rules that satisfy the condition.",
			},
			"lb_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of LB Rules. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The path of rule.",
						},
						"host_header_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the host header.",
						},
						"backend_server_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the backend server group.",
						},
						"listener_sync": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to synchronizethe the health check, the session hold and the forward algorithms of the listener.Valid Values:'on', 'off'.",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Forwarding mode of listener.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation for LB Rule.",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the rule.",
						},

						"health_check": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "A list of health checks.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"health_check_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the health check.",
									},
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the listener.",
									},
									"health_check_state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status maintained by health examination.",
									},
									"healthy_threshold": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Health threshold.",
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
										Description: "Domain name of HTTP type health check.",
									},
								},
							},
							Computed: true,
						},
						"session": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Computed:    true,
							Description: "session configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"session_persistence_period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Session hold timeout.",
									},
									"session_state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The state of session.",
									},
									"cookie_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the cookie.",
									},
									"cookie_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of cookie.",
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

func dataSourceKsyunSlbRulesRead(d *schema.ResourceData, meta interface{}) error {
	slbService := SlbService{meta.(*KsyunClient)}
	return slbService.ReadAndSetLbRules(d, dataSourceKsyunSlbRules())
}
