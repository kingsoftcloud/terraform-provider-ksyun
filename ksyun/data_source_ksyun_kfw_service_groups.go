/*
This data source provides a list of Cloud Firewall Service Group resources according to their instance ID, service group ID, name, and other filters.

# Example Usage

```hcl

	data "ksyun_kfw_service_groups" "default" {
	  output_file = "output_result"
	  cfw_instance_id = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
	  ids = []
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKfwServiceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKfwServiceGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Service Group IDs.",
			},
			"cfw_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud Firewall Instance ID.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Service Groups that satisfy the condition.",
			},
			"kfw_service_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service Group ID.",
						},
						"cfw_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud Firewall Instance ID.",
						},
						"service_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service group name.",
						},
						"service_infos": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Service information. Format: protocol:src_port_min-src_port_max/dest_port_min-dest_port_max. Example: TCP:1-100/70-80, UDP:22/33, ICMP.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description.",
						},
						"citation_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of references.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKfwServiceGroupsRead(d *schema.ResourceData, meta interface{}) error {
	kfwService := KfwService{meta.(*KsyunClient)}
	return kfwService.ReadAndSetKfwServiceGroups(d, dataSourceKsyunKfwServiceGroups())
}
