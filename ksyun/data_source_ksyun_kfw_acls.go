/*
This data source provides a list of Cloud Firewall ACL Rule resources according to their instance ID, ACL ID, name, and other filters.

# Example Usage

```hcl

	data "ksyun_kfw_acls" "default" {
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

func dataSourceKsyunKfwAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKfwAclsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of ACL Rule IDs.",
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
				Description: "Total number of ACL Rules that satisfy the condition.",
			},
			"kfw_acls": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"acl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ACL Rule ID.",
						},
						"cfw_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud Firewall Instance ID.",
						},
						"acl_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ACL rule name.",
						},
						"direction": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Direction. Valid values: in (inbound), out (outbound).",
						},
						"src_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source address type. Valid values: ip, addrbook, zone, any.",
						},
						"src_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Source IP addresses.",
						},
						"src_addrbooks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Source address book IDs.",
						},
						"src_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeMap,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							Description: "Source zones (geographic regions).",
						},
						"dest_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination address type. Valid values: ip, addrbook, any.",
						},
						"dest_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Destination IP addresses.",
						},
						"dest_addrbooks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Destination address book IDs.",
						},
						"dest_host": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Destination domain names.",
						},
						"dest_hostbook": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Destination host book IDs.",
						},
						"service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service type. Valid values: service, servicegroup, any.",
						},
						"service_infos": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Service information. Format: protocol:src_port_min-src_port_max/dest_port_min-dest_port_max. Example: TCP:1-100/80-80, UDP:22/33, ICMP.",
						},
						"service_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Service group IDs.",
						},
						"app_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application type. Valid values: app, any.",
						},
						"app_value": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Application values.",
						},
						"policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action. Valid values: accept, deny.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status. Valid values: start, stop.",
						},
						"priority_position": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Priority position. Format: after+priority or before+priority. Example: after+1, before+1.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Priority value.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description.",
						},
						"hit_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Hit count.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKfwAclsRead(d *schema.ResourceData, meta interface{}) error {
	kfwService := KfwService{meta.(*KsyunClient)}
	return kfwService.ReadAndSetKfwAcls(d, dataSourceKsyunKfwAcls())
}
