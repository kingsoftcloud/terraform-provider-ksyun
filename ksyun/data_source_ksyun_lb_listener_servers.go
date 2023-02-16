/*
Provides a Load Balancer Listener server resource.

# Example Usage

```hcl

	resource "ksyun_lb_listener_server" "default" {
	  listener_id = "3a520244-ddc1-41c8-9d2b-66b4cf3a2386"
	  real_server_ip = "10.0.77.20"
	  real_server_port = 8000
	  real_server_type = "host"
	  instance_id = "3a520244-ddc1-41c8-9d2b-66b4cf3a2386"
	  weight = 10
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunLbListenerServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunLbServersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of LB Listener Server IDs, all the LB Listener Servers belong to this region will be retrieved if the ID is `\"\"`.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of LB Listener Servers that satisfy the condition.",
			},
			"listener_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of LB Listener IDs.",
			},
			"real_server_ip": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of real servers.",
			},
			"servers": {
				Type:        schema.TypeList,
				Computed:    true,
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
						"real_server_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the real server.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Weight of the real server.",
						},
						"real_server_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State of the real server.",
						},

						"register_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "register ID of the real server.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID of the real server, if real server type is host.",
						},
						"master_slave_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "whether real server is master or salve. when listener method is MasterSlave, this field is supported.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunLbServersRead(d *schema.ResourceData, meta interface{}) error {
	slbService := SlbService{meta.(*KsyunClient)}
	return slbService.ReadAndSetRealServers(d, dataSourceKsyunLbListenerServers())
}
