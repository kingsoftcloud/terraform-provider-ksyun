/*
This data source provides a list of ALB listener resources according to their ID.

# Example Usage

```hcl

	data "ksyun_alb_listeners" "default" {
		output_file="output_result"
		ids=[]
		alb_id=[]
		acl_id=[]
		protocol=[]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunAlbListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunAlbListenersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of ALB Listener IDs, all the ALB Listeners belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"alb_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more ALB IDs.",
			},
			"acl_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more ACL ID.",
			},
			"protocol": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more Listener protocol.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of ALB listeners that satisfy the condition.",
			},
			"listeners": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of ALB Listeners. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the ALB Listener.",
						},
						"alb_listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the listener.",
						},
						"alb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the ALB Listener.",
						},
						"alb_listener_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the listener.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time.",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Forwarding mode of listener.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol of listener.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The protocol port of listener.",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of certificate.",
						},
						"tls_cipher_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TLS cipher policy.",
						},
						"alb_listener_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of listener.",
						},
						"redirect_alb_listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the redirect ALB listener.",
						},
						"enable_http2": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether enable to HTTP2.",
						},
						"http_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backend Protocol.",
						},
						"health_check": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Health check configuration. It is a nested type which documented below.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"health_check_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the healthcheck.",
									},
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the LB listener.",
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
										Description: "The service host name of the health check, which is available only for the HTTP or HTTPS health check.",
									},
								},
							},
							Computed: true,
						},
						"session": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Computed:    true,
							Description: "session.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunAlbListenersRead(d *schema.ResourceData, meta interface{}) error {
	s := AlbListenerService{meta.(*KsyunClient)}
	return s.ReadAndSetAlbListeners(d, dataSourceKsyunAlbListeners())
}
