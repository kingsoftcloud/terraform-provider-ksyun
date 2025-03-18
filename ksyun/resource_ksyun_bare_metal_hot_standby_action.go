/*
Provides bare metal use hot standby action.

# Example Usage

```hcl
resource "ksyun_bare_metal_hot_standby_action" "foo" {
  host_id = "epc_id"
  hot_standby {
    hot_stand_by_host_id = "hot_standby_id"
    retain_instance_info = "Notretain"
  }
}

```

# Import

Dont Allow to Import
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunBareMetalHotStandbyAction() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunBareMetalHotStandbyActionCreate,
		Read:   resourceKsyunBareMetalHotStandbyActionRead,
		Update: nil,
		Delete: resourceKsyunBareMetalHotStandbyActionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of epc for hot standby.",
			},
			"hot_standby": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hot_stand_by_host_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The id of hot standby.",
						},
						"retain_instance_info": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Whether retain the instance info. Valid Values: `RetainPrivateIP` `Notretain`.",
						},
					},
				},
				Description: "Indicate the hot standby to instead the master Host.",
			},
		},
	}
}

func resourceKsyunBareMetalHotStandbyActionCreate(d *schema.ResourceData, meta interface{}) (err error) {
	bmSrv := BareMetalService{meta.(*KsyunClient)}
	d.SetId(d.Get("host_id").(string))

	call, err := bmSrv.UseHotStandbyCall(d, resourceKsyunBareMetalHotStandbyAction())
	if err != nil {
		return fmt.Errorf("error on creating request call when using hot standby for epc %q, %s", d.Id(), err)
	}

	if call.executeCall != nil {
		err = ksyunApiCallNew([]ApiCall{call}, d, bmSrv.client, true)
		if err != nil {
			return fmt.Errorf("error on using hot standby for epc %q, %s", d.Id(), err)
		}
	}
	return resourceKsyunBareMetalHotStandbyActionRead(d, meta)
}

func resourceKsyunBareMetalHotStandbyActionRead(d *schema.ResourceData, meta interface{}) (err error) {
	d.SetId(d.Get("host_id").(string))
	return nil
}

func resourceKsyunBareMetalHotStandbyActionUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	return nil
}

func resourceKsyunBareMetalHotStandbyActionDelete(d *schema.ResourceData, meta interface{}) (err error) {
	return nil
}
