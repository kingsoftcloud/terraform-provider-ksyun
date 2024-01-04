/*
Provides a Kcrs Repository Instance resource.

Example Usage

```hcl
# Create a Kcrs Repository Instance
resource "ksyun_kcrs_instance" "foo" {
	instance_name = "tfunittest"
	instance_type = "basic"
}


# Create a Kcrs Repository Instance and open public access
resource "ksyun_kcrs_instance" "foo" {
	instance_name = "tfunittest"
	instance_type = "basic"
	open_public_operation = true


	# open public access with external policy that permits an address, ip or cidr, to access this repository
	external_policy {
		entry = "192.168.2.133"
		desc = "ddd"
	}
	external_policy {
		entry = "192.168.2.123/32"
		desc = "ddd"
	}
}
```

Import

KcrsInstance can be imported using the id, e.g.

```
$ terraform import ksyun_kcrs_instance.foo 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
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
				Description: "Repository instance name.",
			},
			"charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "HourlyInstantSettlement",
				ValidateFunc: validation.StringInSlice([]string{"HourlyInstantSettlement"},
					false),
				ForceNew:    true,
				Description: "Charge type of the instance. Valid Values: 'HourlyInstantSettlement'. Default: 'HourlyInstantSettlement'.",
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"basic", "premium"},
					false),
				ForceNew:    true,
				Description: "The type of instance. Valid Values: 'basic', 'premium'.",
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
							Required:    true,
							Description: "External policy entry. Submit to CIDR or IP.",
						},
						"desc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The external policy description.",
						},
					},
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if !d.Get("open_public_operation").(bool) {
						return true
					}
					return false
				},
				Description: "The external access policy. It's activated when 'open_public_operation' is true.",
			},

			"delete_bucket": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,

				// never activate diff
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Description: "Whether delete bucket with this instance is removing.",
			},

			// computed values
			"instance_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Repository instance status.",
			},
			"internal_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal endpoint address.",
			},
			"public_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public domain.",
			},
			"expired_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expired time.",
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
