/*
This data source provides a list of healthcheck resources  according to their healthcheck ID or listener ID.

# Example Usage

```hcl

	data "ksyun_health_checks" "default" {
	  output_file="output_result"
	  ids=[]
	  listener_id=["8d1dac22-6c6c-42ea-93e2-2702d44ddb93","70467f7e-23dc-465a-a609-fb1525fc6b16"]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunHealthChecks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunHealthChecksRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of health check IDs, all the healthcheck belong to this region will be retrieved if the ID is `\"\"`.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of healthcheck that satisfy the condition.",
			},
			"listener_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of listener IDs, all the healthcheck belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"health_checks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Interval of health examination.",
						},
						"unhealthy_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unhealthy threshold.",
						},
						"health_check_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status maintained by health examination.",
						},
						"health_check_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the healthcheck.",
						},
						"healthy_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Health threshold.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Health check timeout.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the listener.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunHealthChecksRead(d *schema.ResourceData, meta interface{}) error {
	slbService := SlbService{meta.(*KsyunClient)}
	return slbService.ReadAndSetHealthChecks(d, dataSourceKsyunHealthChecks())
}
