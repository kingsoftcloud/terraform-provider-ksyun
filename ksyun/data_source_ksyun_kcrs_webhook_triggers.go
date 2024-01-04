/*
This data source provides a list of webhook trigger resources according to their instance id.

Example Usage

```hcl
data "ksyun_kcrs_webhook_triggers" "foo" {
  output_file="kcrs_webhook_triggers_output_result"
  instance_id = "86f14f8c-bf24-42c8-91bd-xxxxxxxx"
  namespace = "tftest"
}
```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKcrsWebhookTriggers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKcrsWebhookTriggersRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,

				Description: "Kcrs Instance Id.",
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,

				Description: "Kcrs Instance Namespace.",
			},

			"trigger_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Webhook Trigger ID, all the Webhook Trigger belong to this namespace of instance will be retrieved if the ID is `\"\"`.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of kcrs webhook_triggers that satisfy the condition.",
			},

			"triggers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the certificate.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "name of the certificate.",
						},
						"trigger_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the certificate.",
						},
						"trigger_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the certificate.",
						},
						"event_type": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "name of the certificate.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the certificate.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the certificate.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKcrsWebhookTriggersRead(d *schema.ResourceData, meta interface{}) error {
	kcmService := KcrsService{meta.(*KsyunClient)}
	return kcmService.ReadAndSetKcrsWebhookTriggers(d, dataSourceKsyunKcrsWebhookTriggers())
}
