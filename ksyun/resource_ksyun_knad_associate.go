package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunKnadAssociate() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKnadAssociateCreate,
		Read:   resourceKsyunKnadAssociateRead,
		Delete: resourceKsyunKnadAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: importKnadAssociate,
		},

		Schema: map[string]*schema.Schema{
			"knad_id": {
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"ip": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceKsyunKnadAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	knadService := KnadService{meta.(*KsyunClient)}
	err = knadService.AssociateKnad(d, resourceKsyunKnadAssociate())
	if err != nil {
		return fmt.Errorf("error on associate bandWidthShare %q, %s", d.Id(), err)
	}
	return resourceKsyunKnadAssociateRead(d, meta)
}

func resourceKsyunKnadAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	knadService := KnadService{meta.(*KsyunClient)}
	err = knadService.ReadAndSetAssociateKnad(d, resourceKsyunKnadAssociate())
	if err != nil {
		return fmt.Errorf("error on reading bandWidthShare associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunKnadAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	knadService := KnadService{meta.(*KsyunClient)}
	err = knadService.DisassociateKnad(d)
	if err != nil {
		return fmt.Errorf("error on disAssociate bandWidthShare %q, %s", d.Id(), err)
	}
	return err
}
