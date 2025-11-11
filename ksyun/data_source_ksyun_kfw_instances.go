/*
This data source provides a list of Cloud Firewall Instance resources according to their instance ID, name, and other filters.

# Example Usage

```hcl

	data "ksyun_kfw_instances" "default" {
	  output_file = "output_result"
	  ids = []
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKfwInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKfwInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Cloud Firewall Instance IDs.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Cloud Firewall Instances that satisfy the condition.",
			},

			"kfw_instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cfw_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Cloud Firewall Instance.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Cloud Firewall Instance.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type. Valid values: Advanced, Enterprise.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bandwidth (10-5000M).",
						},
						"total_eip_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of protected IPs.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing type. Valid values: Monthly (prepaid), Daily (pay-as-you-go, trial).",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project ID.",
						},
						"purchase_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Purchase duration.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status (1-creating, 2-running, 3-modifying, 4-stopped, 5-abnormal, 6-unsubscribing).",
						},
						"used_eip_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of protected IPs in use.",
						},
						"total_acl_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of ACL rules that can be added.",
						},
						"ips_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IPS status (0-stopped, 1-enabled).",
						},
						"av_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "AV status (0-stopped, 1-enabled).",
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

func dataSourceKsyunKfwInstancesRead(d *schema.ResourceData, meta interface{}) error {
	kfwService := KfwService{meta.(*KsyunClient)}
	return kfwService.ReadAndSetKfwInstances(d, dataSourceKsyunKfwInstances())
}
