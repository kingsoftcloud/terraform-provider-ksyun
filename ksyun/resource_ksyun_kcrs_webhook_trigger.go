/*
Provides a webhook trigger under kcrs repository instance.

Example Usage

```hcl
# repository instance
resource "ksyun_kcrs_instance" "foo" {
	instance_name = "tfunittest"
	instance_type = "basic"
}

# Create a webhook trigger
resource "ksyun_kcrs_webhook_trigger" "foo" {
	instance_id = ksyun_kcrs_instance.foo.id
	namespace = "namespace"
	trigger {
		trigger_url = "http://www.test111.com"
		trigger_name = "tfunittest"
		event_types = ["DeleteImage", "PushImage"]
		headers {
			key = "pp1"
			value = "22"
		}
		headers {
			key = "pp1"
			value = "333"
		}
	}
}
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunKcrsWebhookTrigger() *schema.Resource {
	return &schema.Resource{
		Create:   resourceKsyunKcrsWebhookTriggerCreate,
		Read:     resourceKsyunKcrsWebhookTriggerRead,
		Update:   resourceKsyunKcrsWebhookTriggerUpdate,
		Delete:   resourceKsyunKcrsWebhookTriggerDelete,
		Importer: &schema.ResourceImporter{
			// State: importKcrsWebhookTrigger,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     0,
				Description: "Instance id of repository.",
			},

			"namespace": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Namespace name.",
			},

			"trigger": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "trigger parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger_url": {
							Required:    true,
							Type:        schema.TypeString,
							Description: "The post url for webhook after trigger action.",
						},
						"headers": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Custom Headers.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Header Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Header Value.",
									},
								},
							},
						},
						"event_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"PushImage", "DeleteImage"}, false),
							},
							Required:    true,
							Description: "Trigger action. Valid Values: 'PushImage', 'DeleteImage'.",
						},
						"trigger_name": {
							Type:     schema.TypeString,
							Required: true,
							// ForceNew:    true,
							Description: "Trigger name.",
						},

						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enable trigger.",
						},
					},
				},
			},
		},
	}
}
func resourceKsyunKcrsWebhookTriggerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kcrsWebhookTriggerService := KcrsService{meta.(*KsyunClient)}
	err = kcrsWebhookTriggerService.CreateKcrsWebhookTrigger(d, resourceKsyunKcrsWebhookTrigger())
	if err != nil {
		return fmt.Errorf("error on creating kcrs WebhookTrigger %q, %s", d.Id(), err)
	}
	return resourceKsyunKcrsWebhookTriggerRead(d, meta)
}

func resourceKsyunKcrsWebhookTriggerRead(d *schema.ResourceData, meta interface{}) (err error) {
	kcrsWebhookTriggerService := KcrsService{meta.(*KsyunClient)}
	err = kcrsWebhookTriggerService.ReadAndSetWebhookTrigger(d, resourceKsyunKcrsWebhookTrigger())
	if err != nil {
		return fmt.Errorf("error on reading kcrs WebhookTrigger %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunKcrsWebhookTriggerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	kcrsWebhookTriggerService := KcrsService{meta.(*KsyunClient)}
	err = kcrsWebhookTriggerService.ModifyWebhookTrigger(d, resourceKsyunKcrsWebhookTrigger())
	if err != nil {
		return fmt.Errorf("error on updating kcrs WebhookTrigger %q, %s", d.Id(), err)
	}
	return resourceKsyunKcrsWebhookTriggerRead(d, meta)
}

func resourceKsyunKcrsWebhookTriggerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	kcrsWebhookTriggerService := KcrsService{meta.(*KsyunClient)}
	err = kcrsWebhookTriggerService.RemoveKcrsWebhookTrigger(d)
	if err != nil {
		return fmt.Errorf("error on deleting kcrs WebhookTrigger %q, %s", d.Id(), err)
	}
	return err

}
