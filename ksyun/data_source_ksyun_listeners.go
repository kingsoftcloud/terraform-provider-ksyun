/*
This data source provides a list of Load Balancer Listener resources according to their Load Balancer Listener ID.

# Example Usage

```hcl

	data "ksyun_listeners" "default" {
	  output_file="output_result"
	  ids=[""]
	  load_balancer_id=["d3fd0421-a35a-4ddb-a939-5c51e8af8e8c","4534d617-9de0-4a4a-9ed5-3561196cacb6"]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunListenersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of LB Listener IDs, all the LB Listeners belong to this region will be retrieved if the ID is `\"\"`.",
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter resulting lb listeners by name.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of LB listeners that satisfy the condition.",
			},
			"load_balancer_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of load balancer IDs.",
			},
			"certificate_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of certificate IDs.",
			},
			"listeners": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation.",
						},
						"load_balancer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the LB.",
						},
						"listener_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State of the LB listener.",
						},
						"listener_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the LB listener.",
						},
						"listener_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol of the LB listener.",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the certificate.",
						},
						"listener_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port of the LB listener.",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The load balancer method in which the listener is.",
						},

						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the LB listener.",
						},
						"enable_http2": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether support HTTP2.",
						},

						"tls_cipher_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Https listener TLS cipher policy.",
						},
						"http_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HTTP protocol.",
						},

						"load_balancer_acl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the LB ACL ID.",
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
						"real_server": {
							Type:        schema.TypeList,
							Description: "An information list of real servers. Each element contains the following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"real_server_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP of the real server.",
									},
									"real_server_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port of the real server.",
									},
									"real_server_state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "State of the real server.",
									},
									"real_server_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the real server.",
									},
									"register_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "register ID of the real server.",
									},
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the LB listener.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ID of the real server, if real server type is host.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Weight of the real server.",
									},
								},
							},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunListenersRead(d *schema.ResourceData, meta interface{}) error {
	slbService := SlbService{meta.(*KsyunClient)}
	return slbService.ReadAndSetListeners(d, dataSourceKsyunListeners())
}
