/*
Provides a list of lb host headers in the current region.

# Example Usage

```hcl

	data "ksyun_lb_host_headers" "default" {
		output_file="output_result"
		ids=[]
		listener_id=[]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunListenerHostHeaders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunListenerHostHeadersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of hostheader IDs.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running terraform plan).",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of host headers that satisfy the condition.",
			},
			"listener_id": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of the listeners.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"host_headers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of host headers. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the listener.",
						},
						"host_header": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the host header.",
						},
						"host_header_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the host header.",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of certificate, HTTPS type listener creates this parameter which is not default.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunListenerHostHeadersRead(d *schema.ResourceData, meta interface{}) error {
	slbService := SlbService{meta.(*KsyunClient)}
	return slbService.ReadAndSetHostHeaders(d, dataSourceKsyunListenerHostHeaders())
}
