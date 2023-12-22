/*
Provides an KcrsInstance resource.

Example Usage

```hcl
# Create a KcrsInstance
resource "ksyun_KcrsInstance" "default" {
  link_type = "DDoS_BGP"
  ip_count = 10
  band = 30
  max_band = 30
  idc_band = 100
  duration = 1
  KcrsInstance_name = "ksc_KcrsService"
  bill_type = 1
  service_id = "KcrsInstance_30G"
  project_id="0"
}
```

Import

KcrsInstance can be imported using the id, e.g.

```
$ terraform import ksyun_KcrsInstance.default KcrsService67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunKcrsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKcrsInstanceCreate,
		Read:   resourceKsyunKcrsInstanceRead,
		Update: resourceKsyunKcrsInstanceUpdate,
		Delete: resourceKsyunKcrsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     0,
				Description: "The id of the project.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the ID of the KcrsInstance.",
			},
			"charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "HourlyInstantSettlement",
				ValidateFunc: validation.StringInSlice([]string{"HourlyInstantSettlement"},
					false),
				ForceNew:    true,
				Description: "the link type of the KcrsInstance. Valid Values: 'HourlyInstantSettlement'.",
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"basic", "premium"},
					false),
				ForceNew:    true,
				Description: "the max ip count that can bind to the KcrsInstance,value range: [10, 100].",
			},

			"open_public_operation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Control public network access.",
			},

			"external_policy": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry": {
							Type: schema.TypeString,
							ValidateFunc: validation.Any(
								validation.IsCIDR,
								validation.IsIPAddress,
							),
							Required: true,
						},
						"desc": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if !d.Get("open_public_operation").(bool) {
						return true
					}
					return false
				},
			},

			"delete_bucket": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,

				// never activate diff
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
			},
			// computed values
			"instance_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"internal_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"public_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"expired_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
		},
	}
}
func resourceKsyunKcrsInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	KcrsInstanceService := KcrsService{meta.(*KsyunClient)}
	err = KcrsInstanceService.CreateKcrsInstance(d, resourceKsyunKcrsInstance())
	if err != nil {
		return fmt.Errorf("error on creating kcrs instance %q, %s", d.Id(), err)
	}
	return resourceKsyunKcrsInstanceRead(d, meta)
}

func resourceKsyunKcrsInstanceRead(d *schema.ResourceData, meta interface{}) (err error) {

	KcrsInstanceService := KcrsService{meta.(*KsyunClient)}
	err = KcrsInstanceService.ReadAndSetKcrsInstance(d, resourceKsyunKcrsInstance())
	if err != nil {
		return fmt.Errorf("error on reading kcrs instance %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunKcrsInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	kcrsInstanceService := KcrsService{meta.(*KsyunClient)}

	if d.HasChanges("open_public_operation", "external_policy") {
		err := kcrsInstanceService.ModifyKcrsInstanceEoIEndpoint(d, resourceKsyunKcrsInstance())
		if err != nil {
			return fmt.Errorf("an error caused when changing instance external endpoint status or policy %q, %s", d.Id(), err)
		}
		return resourceKsyunKcrsInstanceRead(d, meta)
	}

	return fmt.Errorf("the attributes, 'project_id', are not supported to modify")
}

func resourceKsyunKcrsInstanceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	KcrsInstanceService := KcrsService{meta.(*KsyunClient)}
	err = KcrsInstanceService.RemoveKcrsInstance(d)
	if err != nil {
		return fmt.Errorf("error on deleting kcrs instance %q, %s", d.Id(), err)
	}
	return err
}
