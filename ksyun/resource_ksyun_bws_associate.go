/*
Provides a BWS Association resource for associating EIP with a BWS instance.

# Example Usage

```hcl

	resource "ksyun_bws_associate" "default" {
	  band_width_share_id = "2af77683-b47e-4634-88ce-fcb95cb65e86"
	  allocation_id = "139134fc-f622-467f-a8b1-c0858dac62ab"
	}

```

# Import

# BWS can be imported using the id

```
$ terraform import ksyun_bws_associate.default ${band_width_share_id}:${allocation_id}
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunBandWidthShareAssociate() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunBandWidthShareAssociateCreate,
		Read:   resourceKsyunBandWidthShareAssociateRead,
		Delete: resourceKsyunBandWidthShareAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: importBandWidthShareAssociate,
		},

		Schema: map[string]*schema.Schema{
			"band_width_share_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the BWS.",
			},
			"allocation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the EIP.",
			},
			"band_width": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "bandwidth value.",
			},
		},
	}
}

func resourceKsyunBandWidthShareAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.AssociateBandWidthShare(d, resourceKsyunBandWidthShareAssociate())
	if err != nil {
		return fmt.Errorf("error on associate bandWidthShare %q, %s", d.Id(), err)
	}
	return resourceKsyunBandWidthShareAssociateRead(d, meta)
}

func resourceKsyunBandWidthShareAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.ReadAndSetAssociateBandWidthShare(d, resourceKsyunBandWidthShareAssociate())
	if err != nil {
		return fmt.Errorf("error on reading bandWidthShare associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunBandWidthShareAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.DisassociateBandWidthShare(d)
	if err != nil {
		return fmt.Errorf("error on disAssociate bandWidthShare %q, %s", d.Id(), err)
	}
	return err

}
