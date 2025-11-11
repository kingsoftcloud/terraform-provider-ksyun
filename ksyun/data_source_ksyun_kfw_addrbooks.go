/*
This data source provides a list of Cloud Firewall Address Book resources according to their instance ID, address book ID, name, and other filters.

# Example Usage

```hcl

	data "ksyun_kfw_addrbooks" "default" {
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

func dataSourceKsyunKfwAddrbooks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKfwAddrbooksRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Address Book IDs.",
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
				Description: "Total number of Address Books that satisfy the condition.",
			},
			"kfw_addrbooks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"addrbook_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Address Book ID.",
						},
						"cfw_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud Firewall Instance ID.",
						},
						"addrbook_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Address book name.",
						},
						"ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP version. Valid values: IPv4, IPv6.",
						},
						"ip_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "IP addresses.",
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

func dataSourceKsyunKfwAddrbooksRead(d *schema.ResourceData, meta interface{}) error {
	kfwService := KfwService{meta.(*KsyunClient)}
	return kfwService.ReadAndSetKfwAddrbooks(d, dataSourceKsyunKfwAddrbooks())
}
