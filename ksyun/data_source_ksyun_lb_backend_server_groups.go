/*
Provides a list of lb backend server groups in the current region.

# Example Usage

```hcl
# Get availability zones
data "ksyun_lb_backend_server_groups" "default" {
output_file="out_file"
ids=[]
}
```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunBackendServerGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunBackendServerGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of BackendServerGroup IDs.",
			},
			"vpc_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of VPC IDs.",
			},
			"backend_server_group_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of BackendServerGroup types.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of BackendServerGroups that satisfy the condition.",
			},
			"backend_server_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of BackendServerGroups. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend_server_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of backend server group.",
						},
						"backend_server_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of backend server group.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual private network ID.",
						},
						"backend_server_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of backend server number.",
						},
						"backend_server_group_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of backend server group.Valid values are Server and Mirror.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol of the backend server group. Valid values: 'TCP', 'UDP', 'HTTP'. Default `HTTP`.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when the backend server group was created.",
						},
						"health_check": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Health check information, only the mirror server has this parameter.",
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
										Description: "state of the health check.",
									},
									"healthy_threshold": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "health threshold.",
									},
									"interval": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "interval of the health check.",
									},
									"timeout": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "timeout of the health check.",
									},
									"unhealthy_threshold": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Unhealthy threshold of health check.",
									},
									"url_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "path of the health check.",
									},
									"host_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "hostname of the health check.",
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

func dataSourceKsyunBackendServerGroupsRead(d *schema.ResourceData, meta interface{}) error {
	slbService := SlbService{meta.(*KsyunClient)}
	return slbService.ReadAndSetBackendServerGroups(d, dataSourceKsyunBackendServerGroups())
}
