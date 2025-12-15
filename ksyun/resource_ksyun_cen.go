/*
Provides a Cen resource.

# Example Usage

```hcl

	resource "ksyun_cen" "default" {
	  cen_name="cen_create"
	  description="zice_create"
	}

```

# Import

Cen can be imported using the `id`, e.g.

```
$ terraform import ksyun_cen.example xxxxxxxx-abc123456
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunCen() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunCenCreate,
		Update: resourceKsyunCenUpdate,
		Read:   resourceKsyunCenRead,
		Delete: resourceKsyunCenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "creation time of the cen.",
			},
			"cen_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the cen.",
			},
			"cen_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the cen.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cen Description.",
			},
		},
	}
}

func resourceKsyunCenCreate(d *schema.ResourceData, meta interface{}) (err error) {
	cenService := CenService{meta.(*KsyunClient)}
	err = cenService.CreateCen(d, resourceKsyunCen())
	if err != nil {
		return fmt.Errorf("error on creating cen %q, %s", d.Id(), err)
	}
	return resourceKsyunCenRead(d, meta)
}

func resourceKsyunCenRead(d *schema.ResourceData, meta interface{}) (err error) {
	cenService := CenService{meta.(*KsyunClient)}
	err = cenService.ReadAndSetCen(d, resourceKsyunCen())
	if err != nil {
		return fmt.Errorf("error on reading cen %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunCenUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	cenService := CenService{meta.(*KsyunClient)}
	err = cenService.ModifyCen(d, resourceKsyunCen())
	if err != nil {
		return fmt.Errorf("error on updating cen %q, %s", d.Id(), err)
	}
	return resourceKsyunCenRead(d, meta)
}

func resourceKsyunCenDelete(d *schema.ResourceData, meta interface{}) (err error) {
	cenService := CenService{meta.(*KsyunClient)}
	err = cenService.RemoveCen(d)
	if err != nil {
		return fmt.Errorf("error on deleting cen %q, %s", d.Id(), err)
	}
	return err
}
