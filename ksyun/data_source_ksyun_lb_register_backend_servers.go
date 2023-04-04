/*
Provides a list of register backend servers in the current region.

# Example Usage

```hcl

	data "ksyun_lb_register_backend_servers" "foo" {
		output_file="output_result"
		ids=[]
		backend_server_group_id=[]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunRegisterBackendServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunRegisterBackendServersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Register backend server IDs, all the Register backend servers belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"backend_server_group_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Register backend server IDs.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Register backend groups that satisfy the condition.",
			},
			"register_backend_servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of Register backend groups. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend_server_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the server.",
						},
						"backend_server_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP of the server.",
						},
						"weight": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The weight of backend service.",
						},
						"register_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The registration ID of the binding server group.",
						},
						"real_server_ip": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The IP of real server.",
						},
						"real_server_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The port number of real server.",
						},
						"real_server_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of real server.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the instance.",
						},
						"network_interface_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the network interface.",
						},
						"real_server_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the real server.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunRegisterBackendServersRead(d *schema.ResourceData, meta interface{}) error {
	slbService := SlbService{meta.(*KsyunClient)}
	return slbService.ReadAndSetBackendServers(d, dataSourceKsyunRegisterBackendServers())
}
