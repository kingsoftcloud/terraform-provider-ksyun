/*
Provides a PerPay KNAD (PerKnad) resource.

# Example Usage

```hcl
# Create a perpay knad

	resource "ksyun_perknad" "default" {
	  ip_count   = 10
	  max_band   = 30
	  knad_name  = "ksc_kad"
	  project_id = "0"
	}

```

# Import

PerKnad can be imported using the id, e.g.

```
$ terraform import ksyun_perknad.default knad67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunPerKnad() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunPerKnadCreate,
		Read:   resourceKsyunPerKnadRead,
		Update: resourceKsyunPerKnadUpdate,
		Delete: resourceKsyunPerKnadDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0",
				Description: "The id of the project.",
			},
			"knad_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the ID of the PerKnad.",
			},
			"knad_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the name of the PerKnad.",
			},
			"ip_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(10, 100),
				Description:  "the max ip count that can bind to the PerKnad, value range: [10, 100].",
			},
			"max_band": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "the max protection band of the PerKnad.",
			},
		},
	}
}

func resourceKsyunPerKnadCreate(d *schema.ResourceData, meta interface{}) (err error) {
	svc := PerKnadService{meta.(*KsyunClient)}
	err = svc.CreatePerKnad(d, resourceKsyunPerKnad())
	if err != nil {
		return fmt.Errorf("error on creating perknad %q, %s", d.Id(), err)
	}
	return resourceKsyunPerKnadRead(d, meta)
}

func resourceKsyunPerKnadRead(d *schema.ResourceData, meta interface{}) (err error) {
	svc := PerKnadService{meta.(*KsyunClient)}
	err = svc.ReadAndSetPerKnad(d, resourceKsyunPerKnad())
	if err != nil {
		return fmt.Errorf("error on reading perknad %q, %s", d.Id(), err)
	}
	return nil
}

func resourceKsyunPerKnadUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	svc := PerKnadService{meta.(*KsyunClient)}
	err = svc.ModifyPerKnad(d, resourceKsyunPerKnad())
	if err != nil {
		return fmt.Errorf("error on updating perknad %q, %s", d.Id(), err)
	}
	return resourceKsyunPerKnadRead(d, meta)
}

func resourceKsyunPerKnadDelete(d *schema.ResourceData, meta interface{}) (err error) {
	svc := PerKnadService{meta.(*KsyunClient)}
	err = svc.RemovePerKnad(d)
	if err != nil {
		return fmt.Errorf("error on deleting perknad %q, %s", d.Id(), err)
	}
	return nil
}
